package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/bam_agent/tool"
	"github.com/blockassets/cgminer_client"
)

/*
	Describes all the configs for the implemented monitors
*/
type Config struct {
	HighLoad       HighLoadConfig `json:"highLoad"`
	Reboot         RebootConfig   `json:"reboot"`
	CGMQuit        CGMQuitConfig  `json:"cgMinerQuit"`
	AcceptedShares AcceptedConfig `json:"acceptedShares"`
	HighTemp       HighTempConfig `json:"highTemperature"`
}

/*
	Monitors have a simple life cycle interface of just Start/Stop.
*/
type Monitor interface {
	Start() error
	Stop()
}

/*
	Context is the implementation struct for the Monitor interface.
	Each Monitor has a (ctx *Context) which holds the channel for
*/
type Context struct {
	quit      chan bool
	waitGroup *sync.WaitGroup
}

// For tests only! Uses a new WaitGroup, so don't use when a Manager is in use since that uses a shared WG.
func makeContext() *Context {
	return &Context{quit: make(chan bool), waitGroup: &sync.WaitGroup{}}
}

/*
	Common stop method for all Monitors. Just shut down the channel.
*/
func (ctx *Context) Stop() {
	close(ctx.quit)
}

/*
	Interface for the Monitor Managers
*/
type Lifecycle interface {
	StartMonitors()
	StopMonitors()
}

/*
	Implements the Lifecycle interface
*/
type Manager struct {
	Config   *Config
	Client   *cgminer_client.Client
	Monitors *[]Monitor
	sync.WaitGroup
	sync.Mutex
}

/*
	Creates a new Context with the common waitGroup of the manager
*/
func (mgr *Manager) NewContext() *Context {
	return &Context{quit: make(chan bool), waitGroup: &mgr.WaitGroup}
}

/*
	Implementation of the Manager interface for starting monitors.
*/
func (mgr *Manager) StartMonitors() {
	// Blocks until all the monitors are finished. Prevents double start.
	mgr.Wait()

	log.Println("Monitors being started")
	statRetriever := service.NewStatRetriever()
	cgQuitFunc := func() { mgr.Client.Quit() }
	onStallFunc := func() { service.Reboot() }
	onHighTempFunc := func() { service.StopMiner() }

	mgr.Lock()
	defer mgr.Unlock()
	mgr.Monitors = &[]Monitor{
		newLoadMonitor(mgr.NewContext(), &mgr.Config.HighLoad, statRetriever, service.Reboot),
		newPeriodicReboot(mgr.NewContext(), &mgr.Config.Reboot, service.Reboot),
		newPeriodicCGMQuit(mgr.NewContext(), &mgr.Config.CGMQuit, cgQuitFunc),
		newAcceptedMonitor(mgr.NewContext(), &mgr.Config.AcceptedShares, mgr.Client, onStallFunc),
		newHighTempMonitor(mgr.NewContext(), &mgr.Config.HighTemp, mgr.Client, onHighTempFunc),
	}
	for _, monitor := range *mgr.Monitors {
		monitor.Start()
	}
}

/*
	Implementation of the Manager interface for stopped monitors. Prevents double stop
	to prevent a panic()
*/
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

/*
	Simple function for creating a new Ticker go fun() based on the Clock interface
*/
func (ctx *Context) makeTickerFunc(doIt func(), period time.Duration) func() {
	return ctx.makeClockFunc(tool.NewTicker(period), doIt)
}

/*
	Simple function for creating a new Timer go fun() based on the Clock interface
*/
func (ctx *Context) makeTimerFunc(doIt func(), period time.Duration) func() {
	return ctx.makeClockFunc(tool.NewTimer(period), doIt)
}
