package monitor

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type HighLoadConfig struct {
	Enabled         bool    `json:"enabled"`
	PeriodInSeconds int     `json:"periodInSeconds"`
	HighLoadMark    float64 `json:"highLoadMark"`
}

type loadMonitor struct {
	monitorControl // delegate the synchronization and implementation for start, stop etc
	sr             service.StatRetriever
	onHighLoad     func()
}

func newLoadMonitor(sr service.StatRetriever, onHighLoad func()) Monitor {
	return &loadMonitor{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, sr, onHighLoad}
}

func (monitor *loadMonitor) Start(cfgMon *Config) error {
	cfg := cfgMon.Load
	if monitor.IsRunning() {
		return errors.New("loadMonitor:Already started")
	}

	monitor.setRunning()
	monitor.quitter = make(chan struct{})

	go func() {
		log.Printf("Starting Load Monitor: Enabled:%v Checking load > %v every: %v seconds\n", cfg.Enabled, cfg.HighLoadMark, cfg.PeriodInSeconds)
		ticker := time.NewTicker(time.Duration(cfg.PeriodInSeconds) * time.Second)
		defer ticker.Stop()
		defer monitor.stoppedRunning()
		for {
			select {
			case <-ticker.C:
				if cfg.Enabled {
					checkLoad(monitor.sr, cfg.HighLoadMark, monitor.onHighLoad)
				}
			case <-monitor.quitter:
				return
			}
		}
	}()

	return nil
}

func checkLoad(sr service.StatRetriever, highLoadMark float64, onHighLoad func()) (bool, error) {
	loads, err := sr.GetLoad()
	high := false
	if err != nil {
		log.Printf("Error checking LoadAvg: %v", err)
		return high, err
	}
	if loads.FiveMinAvg > highLoadMark {
		high = true
		onHighLoad()
	}
	return high, nil
}
