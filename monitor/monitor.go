package monitor

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/blockassets/bam_agent/service"
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
}

func (mgr *Manager) NewContext() *Context {
	return &Context{quit: make(chan bool), waitGroup: &mgr.WaitGroup}
}

func (mgr *Manager) StartMonitors() {
	// Blocks until all the monitors are finished. Prevents double start.
	mgr.Wait()

	log.Println("Monitors being started")

	periodicRebootInitial := getRandomizedInitialPeriod(mgr.Config.Reboot.Period)
	periodicCGMQuitInitial := getRandomizedInitialPeriod(mgr.Config.CGMQuit.Period)

	statRetriever := service.NewStatRetriever()
	cgQuitFunc := func() { mgr.Client.Quit() }

	mgr.Monitors = &[]Monitor{
		newLoadMonitor(mgr.NewContext(), &mgr.Config.HighLoad, statRetriever, service.Reboot),
		newPeriodicReboot(mgr.NewContext(), &mgr.Config.Reboot, &periodicRebootInitial, service.Reboot),
		newPeriodicCGMQuit(mgr.NewContext(), &mgr.Config.CGMQuit, &periodicCGMQuitInitial, cgQuitFunc),
	}

	for _, monitor := range *mgr.Monitors {
		monitor.Start()
	}
}

func (mgr *Manager) StopMonitors() {

	log.Println("Monitors being stopped")

	for _, monitor := range *mgr.Monitors {
		monitor.Stop()
	}

	// Blocks until all the monitors are finished
	mgr.Wait()
}

/*
	If all miners are reset, they come back on line in a random distribution so that we dont get seen as a
	denial of service attack on the pool. Helper to create randomized initial period
*/
func getRandomizedInitialPeriod(period time.Duration) time.Duration {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return period + time.Duration(r1.Intn(3600))*time.Second
}

/*
	We use a helper function to build the goroutine so that we can encapsulate
	the functionality of having to stop the ticker and inc/dec the waitGroup.
*/
func (ctx *Context) makeTickerFunc(doIt func(), period *time.Duration) func() {
	return func() {
		ctx.waitGroup.Add(1)
		ticker := time.NewTicker(*period)
		for {
			select {
			case <-ticker.C:
				doIt()
			case <-ctx.quit:
				ticker.Stop()
				ctx.waitGroup.Done()
				return
			}
		}
	}
}

/*
	We use a helper function to build the goroutine so that we can encapsulate
	the functionality of having to stop the timer and inc/dec the waitGroup.
*/
func (ctx *Context) makeTimerFunc(doIt func(), period *time.Duration) func() {
	return func() {
		ctx.waitGroup.Add(1)
		timer := time.NewTimer(*period)
		for {
			select {
			case <-timer.C:
				doIt()
			case <-ctx.quit:
				timer.Stop()
				ctx.waitGroup.Done()
				return
			}
		}
	}
}
