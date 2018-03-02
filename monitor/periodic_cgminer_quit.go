package monitor

import (
	"log"
	"time"
)

type CGMQuitConfig struct {
	Enabled         bool `json:"enabled"`
	PeriodInSeconds int  `json:"periodInSeconds"`
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

		go func() {
			monitor.waitGroup.Add(1)
			timer := time.NewTimer(*monitor.initialPeriod)
			for {
				select {
				case <-timer.C:
					monitor.CGMinerQuit()
				case <-monitor.quit:
					timer.Stop()
					monitor.waitGroup.Done()
					return
				}
			}
		}()
	} else {
		log.Println("PeriodicCGMQuitMonitor: Not enabled")
	}

	return nil
}
