package monitor

import (
	"log"
	"time"
)

type RebootConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
}

// Implements the Monitor interface
type PeriodicRebootMonitor struct {
	*Context
	config        *RebootConfig
	initialPeriod time.Duration
	reboot        func()
}

func newPeriodicReboot(context *Context, config *RebootConfig, initialPeriod time.Duration, rebootFunc func()) Monitor {
	return &PeriodicRebootMonitor{
		Context:       context,
		config:        config,
		initialPeriod: initialPeriod,
		reboot:        rebootFunc,
	}
}

func (monitor *PeriodicRebootMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("PeriodicRebootMonitor: reboot in %v", monitor.initialPeriod)

		go monitor.makeTimerFunc(monitor.reboot, monitor.initialPeriod)()
	} else {
		log.Println("PeriodicRebootMonitor: Not enabled")
	}

	return nil
}
