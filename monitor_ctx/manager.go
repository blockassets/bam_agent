package monitor_ctx

import (
	"context"
	"log"
	"sync"

	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/cgminer_client"
)

// Pulls all the monitors together and provides an app level API to start and stop them safely
//
// This utility ensures that the monitors are eitehr all stopped, or are all starrted and there are no orphans
// Starts and stops are counted to ensure monitors are started correctly. See below.
// Start and Stop are synchronized
//
//
type Config struct {
	HighLoad       HighLoadConfig `json:"highLoad"`
	AcceptedShares AcceptedConfig `json:"acceptedShares"`
	HighTemp       HighTempConfig `json:"highTemperature"`
	CGMQuit        CGMQuitConfig  `json:"cgMinerQuit"`
	Reboot         RebootConfig   `json:"reboot"`
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

func NewManager(config *Config, miner *cgminer_client.Client) *Manager {
	mm := &Manager{}
	log.Println("Monitors being started")
	// LoadMonitor dependencies
	sr := service.NewStatRetriever()
	onLoadHigh := func() { service.Reboot() }
	// Accepted share dependencies
	onStallFunc := func() { service.Reboot() }
	// high temp dependencies
	onHighTempFunc := func() { service.StopMiner() }
	// CGMQuit dependencies
	cgmQuitFunc := func() { miner.Quit() }
	// Reboot dependancies
	onRebootFunc := func() { service.Reboot() }

	mm.monitors = &[]Monitor{
		NewLoadMonitor(&config.HighLoad, sr, onLoadHigh),
		NewAcceptedMonitor(&config.AcceptedShares, miner, onStallFunc),
		NewHighTempMonitor(&config.HighTemp, miner, onHighTempFunc),
		NewPeriodicCGMQuit(&config.CGMQuit, cgmQuitFunc),
		NewPeriodicReboot(&config.Reboot, onRebootFunc),
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
// These are essentially the same scenarios...
//
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
