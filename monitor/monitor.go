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
	quitter chan struct{}
	running bool
	mutex   *sync.Mutex
	wg      *sync.WaitGroup
}

//
// iIRunning, setRunning, waitOnRunning and stoppedRunning
// provide synchronization around starting and stopping of the monitor
// there are some tricky edge cases and this ensures only one monitor is running
// for each instance of the specific monitor and that monitor.Stop() blocks until the monitor
// actually ends
// See monitor_load for usage patterns
//
func (mc *monitorControl) IsRunning() bool {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	return mc.running
}

func (mc *monitorControl) waitOnRunning() {
	mc.wg.Wait()
}

func (mc *monitorControl) setRunning() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if mc.running {
		return
	}
	mc.running = true
	mc.wg.Add(1)
}

func (mc *monitorControl) stoppedRunning() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	if !mc.running {
		return
	}
	mc.running = false
	mc.wg.Done()
}

func (mc *monitorControl) Stop() {
	close(mc.quitter)
	mc.waitOnRunning()
}

func StartMonitors(cfg *Config, client *cgminer_client.Client) {
	sr := service.LinuxStatRetriever{}

	log.Println("Monitors being started")

	monitors := []Monitor{
		newLoadMonitor(&sr, service.Reboot),
		newPeriodicReboot(service.Reboot),
		newPeriodicCGMQuit(func() { client.Quit() }),
	}

	for _, monitor := range monitors {
		monitor.Start(cfg)
	}
}
