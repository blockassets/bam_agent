package monitor

import (
	"log"
	"sync"

	"github.com/blockassets/bam_agent/service"
)

type MonitorConfig struct {
	Load LoadConfig `json:"load"`
}

type Monitor interface {
	Start(cfg *MonitorConfig) error
	Stop()
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

func StartMonitors(cfg *MonitorConfig) {
	// Startup the goroutines to do the stuff that needs to be monitored
	sr := service.LinuxStatRetriever{}

	log.Println("Monitors being started")

	mc := newLoadMonitor(&sr, service.Reboot)
	mc.Start(cfg)
}
