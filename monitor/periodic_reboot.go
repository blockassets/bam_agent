package monitor

import (
	"errors"
	"log"
	"sync"
	"time"
)

type RebootConfig struct {
	Enabled                     bool `json:"enabled"`
	PeriodInSeconds             int  `json:"periodInSeconds"`
	InitialPeriodRangeInSeconds int  `json:"initialPeriodRangeInSeconds"`
}

// Implements the Monitor interface
type PeriodicRebootMonitor struct {
	*Context
	reboot func()
}

func newPeriodicReboot(rebootFunc func()) Monitor {
	return &PeriodicRebootMonitor{&Context{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, rebootFunc}
}

func (monitor *PeriodicRebootMonitor) Start(config *Config) error {
	cfg := config.Reboot
	if monitor.IsRunning() {
		return errors.New("periodic Reboot: Already started")
	}

	monitor.StartRunning()
	monitor.quitter = make(chan struct{})

	go func() {
		initialPeriod := getRandomizedInitialPeriod(cfg.PeriodInSeconds, cfg.InitialPeriodRangeInSeconds)
		log.Printf("Starting Periodic Reboot: Enabled: %v reboot in: %v", cfg.Enabled, initialPeriod)
		timer := time.NewTimer(initialPeriod)
		defer monitor.StopRunning()
		for {
			select {
			case <-timer.C:
				if cfg.Enabled {
					monitor.reboot()
				}
			case <-monitor.quitter:
				return
			}
		}
	}()

	return nil
}
