package monitor

import (
	"log"

	"github.com/blockassets/bam_agent/tool"
)

type RebootConfig struct {
	Enabled bool                `json:"enabled"`
	Period  tool.RandomDuration `json:"period"`
}

// Implements the Monitor interface
type PeriodicRebootMonitor struct {
	*Context
	config        *RebootConfig
	reboot        func()
}

func newPeriodicReboot(context *Context, config *RebootConfig, rebootFunc func()) Monitor {
	return &PeriodicRebootMonitor{
		Context:       context,
		config:        config,
		reboot:        rebootFunc,
	}
}

func (monitor *PeriodicRebootMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("PeriodicRebootMonitor: reboot in %v", monitor.config.Period.Duration)

		go monitor.makeTimerFunc(monitor.reboot, monitor.config.Period.Duration)()
	} else {
		log.Println("PeriodicRebootMonitor: Not enabled")
	}

	return nil
}
