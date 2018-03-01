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

type periodicReboot struct {
	monitorControl
	reboot func()
}

func newPeriodicReboot(rebootFunc func()) Monitor {
	return &periodicReboot{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, rebootFunc}
}

func (monitor *periodicReboot) Start(cfgMon *Config) error {
	cfg := cfgMon.Reboot
	if monitor.IsRunning() {
		return errors.New("periodic Reboot: Already started")
	}

	monitor.setRunning()
	monitor.quitter = make(chan struct{})

	go func() {
		initialPeriod := getRandomizedInitialPeriod(cfg.PeriodInSeconds, cfg.InitialPeriodRangeInSeconds)
		log.Printf("Starting Periodic Reboot: Enabled: %v reboot in: %v", cfg.Enabled, initialPeriod)
		timer := time.NewTimer(initialPeriod)
		defer monitor.stoppedRunning()
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
