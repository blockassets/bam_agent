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
	monitors   *[]Monitor
	startCount int
	stopFunc   func()
	sync.Mutex
}

// TODO discuss with Jon about best way to inject dependancies into this function
// Wondering if there should be a struct that combines the config loaded from file
// with the dependancies such as the action functions and StatReceivers and Miner interfaces
// etc...
//
func Init(config *Config, sr service.StatRetriever, onLoadHigh func()) *Manager {
	mm := &Manager{}

	mm.monitors = &[]Monitor{
		NewLoadMonitor(&config.HighLoad, sr, onLoadHigh),
	}

	mm.Start()
	return mm
}

// The usage model allows for out of sequence starts and stops
// I.e. the number of stops should match the number of starts before actually starting
// Example1
// 		Monitors are started via Init					((started)startCount  == 1)
// 		Request A arrives, monitors are stopped....  	((stopped)startCount  == 0)
// 		Request A is being processed
// 		Request B comes in -> Monitors are requested stopped again ((stillStopped)startCount  == -1)
// 		Request B is completed, Monitors are asked to start again  ((not started)startCount  == 0)
//  	Request A is still working...
// 		Request A requests starts monitors again					((started)startCount  == 1)
//
// Example2
// 		Monitors are started via Init						((started)startCount  == 1)
// 		Request A arrives, monitors are stopped.... 		((stopped)startCount  == 0
// 		Request A is being processed
// 		Request B comes in -> Monitors are requested stopped again ((stillStopped)startCount  == -1)
//  	Request A is still working...
// 		Request A requests starts monitors again					((not started)startCount  == 0)
//		Request B is still working
//	 	Request B is completed, Monitors are asked to start again	((started)startCount  == 1)
//
// These are esentially the same scenarios...
//

func (mm *Manager) Start() {
	mm.Lock()
	defer mm.Unlock()
	if mm.startCount == 0 {
		mm.stopFunc = StartMonitors(context.Background(), *mm.monitors)
	}
	mm.startCount++
}

func (mm *Manager) Stop() {
	mm.Lock()
	defer mm.Unlock()
	if mm.startCount == 1 {
		mm.stopFunc()
		mm.stopFunc = nil
	}
	mm.startCount--
}
