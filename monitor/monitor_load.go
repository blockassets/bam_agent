package monitor

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type LoadConfig struct {
	Enabled         bool    `json:"enabled"`
	PeriodInSeconds int     `json:"periodInSeconds"`
	HighLoadMark    float64 `json:"highLoadMark"`
}

type loadMonitor struct {
	monitorControl // delegate the synchronization and implementation for start, stop etc
	sr             service.StatRetriever
	onHighLoad     func()
}

func newLoadMonitor(sr service.StatRetriever, onHighLoad func()) *loadMonitor {
	return &loadMonitor{monitorControl{nil, false, &sync.Mutex{}, &sync.WaitGroup{}}, sr, onHighLoad}
}

func (lm *loadMonitor) Start(cfg *MonitorConfig) error {

	if lm.getRunning() {
		return errors.New("loadMonitor:Already started")
	}
	lm.setRunning()
	lm.quiter = make(chan struct{})
	go func() {
		log.Printf("Starting Load Monitor: Enabled:%v Checking load > %v every: %v seconds\n", cfg.Load.Enabled, cfg.Load.HighLoadMark, cfg.Load.PeriodInSeconds)
		ticker := time.NewTicker(time.Duration(cfg.Load.PeriodInSeconds) * time.Second)
		defer ticker.Stop()
		defer lm.stoppedRunning()
		for {
			select {
			case <-ticker.C:
				if cfg.Load.Enabled {
					checkLoad(lm.sr, cfg.Load.HighLoadMark, lm.onHighLoad)
				}
			case <-lm.quiter:
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
