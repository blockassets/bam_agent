package monitor

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/cgminer_client"
)

type MonitorConfig struct {
	Load    LoadConfig    `json:"load"`
	Reboot  RebootConfig  `json:"reboot"`
	CGMQuit CGMQuitConfig `json:"cgMinerQuit"`
}

type Monitor interface {
	Start(cfg *MonitorConfig) error
	Stop()
}

// If all miners are reset, they come back on line in a random distribution so that we dont get seen as a
// denial of service attack on the pool. Helper to create randomized initial period
func getRandomizedInitialPeriod(periodInSeconds int, rangeInSeconds int) time.Duration {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return time.Duration(periodInSeconds)*time.Second + time.Duration(r1.Intn(rangeInSeconds))*time.Second
}

// Shared functionality to manage starting and stopping and synchronization
// across all the monitors
type monitorControl struct {
	quiter    chan struct{}
	isRunning bool
	mutex     *sync.Mutex
	wg        *sync.WaitGroup
}

//
// getRunning, setRunning, waitOnRunning and stoppedRunning
// provide synchronization around starting and stopping of the monitor
// there are some tricky edge cases and this ensures only one monitor is running
// for each instance of the specific monitor and that monitor.Stop() blocks until the monitor
// actually ends
// See monitor_load for usage patterns
//
func (mc *monitorControl) getRunning() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.isRunning
}

func (mc *monitorControl) waitOnRunning() {
	mc.wg.Wait()
}

func (mc *monitorControl) setRunning() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if mc.isRunning {
		return
	}
	mc.isRunning = true
	mc.wg.Add(1)
	return
}

func (mc *monitorControl) stoppedRunning() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if !mc.isRunning {
		return
	}
	mc.isRunning = false
	mc.wg.Done()
	return
}

func (mc *monitorControl) Stop() {
	close(mc.quiter)
	mc.waitOnRunning()
}

func StartMonitors(cfg *MonitorConfig, client *cgminer_client.Client) {
	// Startup the goroutines to do the stuff that needs to be monitored
	sr := service.LinuxStatRetriever{}

	log.Println("Monitors being started")

	mc := newLoadMonitor(&sr, service.Reboot)
	mc.Start(cfg)

	mr := newPeriodicReboot(service.Reboot)
	mr.Start(cfg)

	// TODO add in how to get access to the cgm_quit functionality
	mcgmQ := newPeriodicCGMQuit(func() { client.Quit() })
	mcgmQ.Start(cfg)
}
