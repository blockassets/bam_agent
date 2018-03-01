package monitor

import (
	"errors"
	"log"
	"sync"
	"time"
)

type CGMQuitConfig struct {
	Enabled                     bool `json:"enabled"`
	PeriodInSeconds             int  `json:"periodInSeconds"`
	InitialPeriodRangeInSeconds int  `json:"initialPeriodRangeInSeconds"`
}

type PeriodicCGMQuitMonitor struct {
	*Context
	CGMinerQuit func()
}

func newPeriodicCGMQuit(CGMQuitFunc func()) Monitor {
	return &PeriodicCGMQuitMonitor{&Context{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, CGMQuitFunc}
}

func (monitor *PeriodicCGMQuitMonitor) Start(config *Config) error {
	cfg := config.CGMQuit
	if monitor.IsRunning() {
		return errors.New("periodic CGMQuit: Already started")
	}

	monitor.StartRunning()
	monitor.quitter = make(chan struct{})

	go func() {
		initialPeriod := getRandomizedInitialPeriod(cfg.PeriodInSeconds, cfg.InitialPeriodRangeInSeconds)
		log.Printf("Starting Periodic CGMQuit: Enabled: %v Initial CGMQuit in: %v, then every %v seconds", cfg.Enabled, initialPeriod, cfg.PeriodInSeconds)
		timer := time.NewTimer(initialPeriod)
		defer monitor.StopRunning()
		for {
			select {
			case <-timer.C:
				timer.Reset(time.Duration(cfg.PeriodInSeconds) * time.Second)
				if cfg.Enabled {
					monitor.CGMinerQuit()
				}
			case <-monitor.quitter:
				return
			}
		}
	}()

	return nil
}
