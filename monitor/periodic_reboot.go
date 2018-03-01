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

func newPeriodicReboot(rebootFunc func()) *periodicReboot {
	return &periodicReboot{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, rebootFunc}
}

func (pr *periodicReboot) Start(cfgMon *MonitorConfig) error {
	cfg := cfgMon.Reboot
	if pr.getRunning() {
		return errors.New("periodic Reboot: Already started")
	}
	pr.setRunning()
	pr.quiter = make(chan struct{})
	go func() {
		initialPeriod := getRandomizedInitialPeriod(cfg.PeriodInSeconds, cfg.InitialPeriodRangeInSeconds)
		log.Printf("Starting Periodic Reboot: Enabled:%v reboot in:%v", cfg.Enabled, initialPeriod)
		timer := time.NewTimer(initialPeriod)
		defer pr.stoppedRunning()
		for {
			select {
			case <-timer.C:
				log.Printf("timer_tick\n")
				if cfg.Enabled {
					log.Printf("timer_tick2\n")
					pr.reboot()
				}
			case <-pr.quiter:
				return
			}
		}
	}()
	return nil
}
