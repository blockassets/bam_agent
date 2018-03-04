package monitor

import (
	"log"
	"time"
)

type CGMQuitConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
}

// Implements the Monitor interface
type PeriodicCGMQuitMonitor struct {
	*Context
	config        *CGMQuitConfig
	initialPeriod *time.Duration
	CGMinerQuit   func()
}

func newPeriodicCGMQuit(context *Context, config *CGMQuitConfig, initialPeriod *time.Duration, CGMQuitFunc func()) Monitor {
	return &PeriodicCGMQuitMonitor{
		Context:       context,
		config:        config,
		initialPeriod: initialPeriod,
		CGMinerQuit:   CGMQuitFunc,
	}
}

func (monitor *PeriodicCGMQuitMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("PeriodicCGMQuitMonitor: cgminer quit in: %v", monitor.initialPeriod)

		go monitor.makeTimerFunc(monitor.CGMinerQuit, monitor.initialPeriod)()
	} else {
		log.Println("PeriodicCGMQuitMonitor: Not enabled")
	}

	return nil
}
