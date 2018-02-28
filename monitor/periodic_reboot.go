package monitor

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

type RebootConfig struct {
	Enabled            bool `json:"enabled"`
	PeriodSecs         int  `json:"period_secs"`
	InitialPeriodRange int  `json:"initial_period_range_secs"`
}

func (cfg *RebootConfig) InitialPeriod() time.Duration {
	// If all miners are reset, they come back on line in a random distribution so that we dont get seen as a
	// denial of service attack on the pool
	s1 := rand.NewSource(time.Now().UnixNano())                                                                   // provide synchronization around starting and stopping of the monitor
	r1 := rand.New(s1)                                                                                            // there are some tricky edge cases and this ensures only one monitor is running
	return time.Duration(cfg.PeriodSecs)*time.Second + time.Duration(r1.Intn(cfg.InitialPeriodRange))*time.Second // for each instance of the periodicReboot and that monitor.Stop() blocks until the monitor
} // actually ends

type periodicReboot struct {
	monitorControl
	reboot func()
}

func newPeriodicReboot(rebootFunc func()) *periodicReboot {
	return &periodicReboot{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, rebootFunc}
}

func (pr *periodicReboot) Start(cfg *MonitorConfig) error {

	if pr.getRunning() {
		return errors.New("Periodic Reboot: Already started")
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
