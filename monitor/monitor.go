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
	Load    HighLoadConfig `json:"load"`
	Reboot  RebootConfig   `json:"reboot"`
	CGMQuit CGMQuitConfig  `json:"cgMinerQuit"`
}

type Monitor interface {
	Start(cfg *Config) error
	Stop()
	IsRunning() bool
	StartRunning()
	StopRunning()
}

// Shared functionality to manage starting and stopping and synchronization
// across all the monitors
type Context struct {
	quitter chan struct{}
	running bool
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

// If all miners are reset, they come back on line in a random distribution so that we dont get seen as a
// denial of service attack on the pool. Helper to create randomized initial period
func getRandomizedInitialPeriod(periodInSeconds int, rangeInSeconds int) time.Duration {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return time.Duration(periodInSeconds)*time.Second + time.Duration(r1.Intn(rangeInSeconds))*time.Second
}


//
// provide synchronization around starting and stopping of the monitor
// there are some tricky edge cases and this ensures only one monitor is running
// for each instance of the specific monitor and that monitor.Stop() blocks until the monitor
// actually ends
// See monitor_load for usage patterns
//
func (ctx *Context) IsRunning() bool {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	return ctx.running
}

func (ctx *Context) waitOnRunning() {
	ctx.wg.Wait()
}

func (ctx *Context) StartRunning() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if ctx.running {
		return
	}
	ctx.running = true
	ctx.wg.Add(1)
}

func (ctx *Context) StopRunning() {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	if !ctx.running {
		return
	}
	ctx.running = false
	ctx.wg.Done()
}

func (ctx *Context) Stop() {
	close(ctx.quitter)
	ctx.waitOnRunning()
}

func StartMonitors(config *Config, client *cgminer_client.Client) {
	statRetriever := service.LinuxStatRetriever{}

	log.Println("Monitors being started")

	monitors := []Monitor{
		newLoadMonitor(&statRetriever, service.Reboot),
		newPeriodicReboot(service.Reboot),
		newPeriodicCGMQuit(func() { client.Quit() }),
	}

	for _, monitor := range monitors {
		monitor.Start(config)
	}
}
