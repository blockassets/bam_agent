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
	initialPeriod *time.Duration
	reboot        func()
}

func newPeriodicReboot(context *Context, config *RebootConfig, initialPeriod *time.Duration, rebootFunc func()) Monitor {
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

		go func() {
			monitor.waitGroup.Add(1)
			timer := time.NewTimer(*monitor.initialPeriod)
			for {
				select {
				case <-timer.C:
					monitor.reboot()
				case <-monitor.quit:
					timer.Stop()
					monitor.waitGroup.Done()
					return
				}
			}
		}()
	} else {
		log.Println("PeriodicRebootMonitor: Not enabled")
	}

	return nil
}
