package monitor

import (
	"log"

	"github.com/blockassets/bam_agent/tool"
)

type CGMQuitConfig struct {
	Enabled bool                `json:"enabled"`
	Period  tool.RandomDuration `json:"period"`
}

// Implements the Monitor interface
type PeriodicCGMQuitMonitor struct {
	*Context
	config      *CGMQuitConfig
	CGMinerQuit func()
}

func newPeriodicCGMQuit(context *Context, config *CGMQuitConfig, CGMQuitFunc func()) Monitor {
	return &PeriodicCGMQuitMonitor{
		Context:     context,
		config:      config,
		CGMinerQuit: CGMQuitFunc,
	}
}

func (monitor *PeriodicCGMQuitMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("PeriodicCGMQuitMonitor: cgminer quit in: %v", monitor.config.Period.Duration)

		go monitor.makeTimerFunc(monitor.CGMinerQuit, monitor.config.Period.Duration)()
	} else {
		log.Println("PeriodicCGMQuitMonitor: Not enabled")
	}

	return nil
}
