package monitor_ctx

import (
	"context"
	"sync"

	"github.com/blockassets/bam_agent/service"
)

// Pulls all the monitors together and provides an app level API to start and stop them safely
//
// Mutliple requests can call start and stop with no guarantee of order
// This utility ensures that the monitors are eitehr all stopped, or are all starrted and there are no orphans
// Access to Start and Stop are synchronized
//
//
type Config struct {
	HighLoad HighLoadConfig `json:"highLoad"`
}

type Manager struct {
	monitors *[]Monitor
	stopFunc func()
	sync.Mutex
}

// TODO discuss with Jon about best way to inject dependancies into this function
// Wondering if there should be a struct that combines the config loaded from file
// with the dependancies such as the action functions and StatReceivers and Miner interfaces
// etc...
//
func Init(config *Config) *Manager {
	mm := &Manager{}
	// Monitor specific dependancies
	onLoadHigh := func() { service.Reboot() }
	sr := service.NewStatRetriever()

	mm.monitors = &[]Monitor{
		NewLoadMonitor(&config.HighLoad, sr, onLoadHigh),
	}

	mm.Start()
	return mm
}

func (mm *Manager) Start() {
	mm.Lock()
	defer mm.Unlock()
	if mm.stopFunc != nil {
		mm.stopFunc()
	}
	mm.stopFunc = StartMonitors(context.Background(), *mm.monitors)
}

func (mm *Manager) Stop() {
	mm.Lock()
	defer mm.Unlock()
	if mm.stopFunc != nil {
		mm.stopFunc()
		mm.stopFunc = nil
	}
}
