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

type periodicCGMQuit struct {
	monitorControl
	CGMinerQuit func()
}

func newPeriodicCGMQuit(CGMQuitFunc func()) *periodicCGMQuit {
	return &periodicCGMQuit{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, CGMQuitFunc}
}

func (pmq *periodicCGMQuit) Start(cfgMon *MonitorConfig) error {
	cfg := cfgMon.CGMQuit
	if pmq.getRunning() {
		return errors.New("periodic CGMQuit: Already started")
	}
	pmq.setRunning()
	pmq.quiter = make(chan struct{})
	go func() {
		initialPeriod := getRandomizedInitialPeriod(cfg.PeriodInSeconds, cfg.InitialPeriodRangeInSeconds)
		log.Printf("Starting Periodic CGMQuit: Enabled:%v Initial CGMQuit in:%v seconds, then every %v seconds", cfg.Enabled, initialPeriod, cfg.PeriodInSeconds)
		timer := time.NewTimer(initialPeriod)
		defer pmq.stoppedRunning()
		for {
			select {
			case <-timer.C:
				timer.Reset(time.Duration(cfg.PeriodInSeconds) * time.Second)
				if cfg.Enabled {
					pmq.CGMinerQuit()
				}
			case <-pmq.quiter:
				return
			}
		}
	}()
	return nil
}
