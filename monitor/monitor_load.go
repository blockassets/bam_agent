package monitor

import (
	"errors"
	"log"
	"time"
)

type LoadConfig struct {
	Enabled      bool    `json:"enabled"`
	PeriodSecs   int     `json:"period_secs"`
	HighLoadMark float64 `json:"high_load_mark"`
}

type loadMonitor struct {
	sr         statRetriever
	ticker     *time.Ticker
	quiter     chan struct{}
	isRunning  bool
	onHighLoad func()
}

func newLoadMonitor(sr statRetriever, onHighLoad func()) *loadMonitor {
	return &loadMonitor{sr, nil, nil, false, onHighLoad}
}

func (lm *loadMonitor) start(cfg *MonitorConfig) error {
	if lm.isRunning {
		return errors.New("loadMonitor:Already started")
	}
	lm.isRunning = true
	lm.quiter = make(chan struct{})
	if cfg.Load.Enabled {
		go func() {
			log.Printf("Starting Load Moniter: Checking load < %v every: %v seconds\n", cfg.Load.HighLoadMark, cfg.Load.PeriodSecs)
			lm.ticker = time.NewTicker(time.Duration(cfg.Load.PeriodSecs) * time.Second)
			defer lm.ticker.Stop()
			defer func() { lm.isRunning = false }()
			for {
				select {
				case <-lm.ticker.C:
					high, err := checkLoadAvg(lm.sr, cfg.Load.HighLoadMark)
					if err != nil {
						log.Printf("Error checking LoadAvg: %v", err)
						return
					}
					if high {
						lm.onHighLoad()
					}
				case <-lm.quiter:
					return
				}
			}
		}()
	}
	return nil
}

func (lm *loadMonitor) stop() {
	close(lm.quiter)
}

func checkLoadAvg(sr statRetriever, highLoadMark float64) (bool, error) {
	loads, err := sr.getLoad()
	if err != nil {
		return false, err
	}
	high := false
	if loads.fiveMinAvg > highLoadMark {
		high = true
	}
	return high, err
}
