package monitor

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type RebootConfig struct {
	Enabled                     bool `json:"enabled"`
	PeriodInSeconds             int  `json:"periodInSeconds"`
	InitialPeriodRangeInSeconds int  `json:"initialPeriodRangeInSeconds"`
}

func (cfg *RebootConfig) InitialPeriod() time.Duration {
	// If all miners are reset, they come back on line in a random distribution so that we dont get seen as a
	// denial of service attack on the pool
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return time.Duration(cfg.PeriodInSeconds)*time.Second + time.Duration(r1.Intn(cfg.InitialPeriodRangeInSeconds))*time.Second
}

type periodicReboot struct {
	monitorControl
	reboot func()
}

func newPeriodicReboot(rebootFunc func()) *periodicReboot {
	return &periodicReboot{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, rebootFunc}
}

func (pr *periodicReboot) Start(cfg *MonitorConfig) error {

	if pr.getRunning() {
		return errors.New("periodic Reboot: Already started")
	}
	pr.setRunning()
	pr.quiter = make(chan struct{})
	go func() {
		initialPeriod := cfg.Reboot.InitialPeriod()
		log.Printf("Starting Periodic Reboot: Enabled:%v reboot in:%v", cfg.Reboot.Enabled, initialPeriod)
		timer := time.NewTimer(initialPeriod)
		defer pr.stoppedRunning()
		for {
			select {
			case <-timer.C:
				if cfg.Reboot.Enabled {
					pr.reboot()
				}
			case <-pr.quiter:
				return
			}
		}
	}()
	return nil
}
