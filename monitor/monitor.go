package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/bam_agent/tool"
	"github.com/blockassets/cgminer_client"
)

type Config struct {
	HighLoad HighLoadConfig `json:"highLoad"`
	Reboot   RebootConfig   `json:"reboot"`
	CGMQuit  CGMQuitConfig  `json:"cgMinerQuit"`
}

type Monitor interface {
	Start() error
	Stop()
}

type Context struct {
	quit      chan bool
	waitGroup *sync.WaitGroup
}

// For tests only! Uses a new WaitGroup, so don't use when a Manager is in use since that uses a shared WG.
func makeContext() *Context {
	return &Context{quit: make(chan bool), waitGroup: &sync.WaitGroup{}}
}

func (ctx *Context) Stop() {
	close(ctx.quit)
}

type Lifecycle interface {
	StartMonitors()
	StopMonitors()
}

// Implements the Lifecycle interface
type Manager struct {
	Config   *Config
	Client   *cgminer_client.Client
	Monitors *[]Monitor
	sync.WaitGroup
	sync.Mutex
}

func (mgr *Manager) NewContext() *Context {
	return &Context{quit: make(chan bool), waitGroup: &mgr.WaitGroup}
}

func (mgr *Manager) StartMonitors() {
	// Blocks until all the monitors are finished. Prevents double start.
	mgr.Wait()

	log.Println("Monitors being started")
	statRetriever := service.NewStatRetriever()
	cgQuitFunc := func() { mgr.Client.Quit() }

	mgr.Lock()
	defer mgr.Unlock()
	mgr.Monitors = &[]Monitor{
		newLoadMonitor(mgr.NewContext(), &mgr.Config.HighLoad, statRetriever, service.Reboot),
		newPeriodicReboot(mgr.NewContext(), &mgr.Config.Reboot, service.Reboot),
		newPeriodicCGMQuit(mgr.NewContext(), &mgr.Config.CGMQuit, cgQuitFunc),
	}
	for _, monitor := range *mgr.Monitors {
		monitor.Start()
	}

}

func (mgr *Manager) StopMonitors() {

	log.Println("Monitors being stopped")
	mgr.Lock()
	defer mgr.Unlock()
	if mgr.Monitors != nil {
		for _, monitor := range *mgr.Monitors {
			monitor.Stop()
		}
		mgr.Monitors = nil
	}

	// Blocks until all the monitors are finished
	mgr.Wait()
}

/*
	We use a helper function to build the goroutine so that we can encapsulate
	the functionality of having to stop the timer/ticker and inc/dec the waitGroup.
*/
func (ctx *Context) makeClockFunc(clock tool.Clock, doIt func()) func() {
	return func() {
		ctx.waitGroup.Add(1)
		for {
			select {
			case <-clock.C():
				doIt()
			case <-ctx.quit:
				clock.Stop()
				ctx.waitGroup.Done()
				return
			}
		}
	}
}

func (ctx *Context) makeTickerFunc(doIt func(), period time.Duration) func() {
	return ctx.makeClockFunc(tool.NewTicker(period), doIt)
}

func (ctx *Context) makeTimerFunc(doIt func(), period time.Duration) func() {
	return ctx.makeClockFunc(tool.NewTimer(period), doIt)
}
