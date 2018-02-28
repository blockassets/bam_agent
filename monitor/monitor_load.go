package monitor

import (
	"errors"
	"log"
	"sync"
	"time"
)

type LoadConfig struct {
	Enabled      bool    `json:"enabled"`
	PeriodSecs   int     `json:"period_secs"`
	HighLoadMark float64 `json:"high_load_mark"`
}

type loadMonitor struct {
	sr         statRetriever
	quiter     chan struct{}
	isRunning  bool
	onHighLoad func()
	mutex      *sync.Mutex
	wg         *sync.WaitGroup
}

func newLoadMonitor(sr statRetriever, onHighLoad func()) *loadMonitor {
	return &loadMonitor{sr, nil, false, onHighLoad, &sync.Mutex{}, &sync.WaitGroup{}}
}

func (lm *loadMonitor) Start(cfg *MonitorConfig) error {

	if lm.getRunning() {
		return errors.New("loadMonitor:Already started")
	}
	lm.setRunning()
	lm.quiter = make(chan struct{})
	go func() {
		log.Printf("Starting Load Moniter: Enabled:%v Checking load < %v every: %v seconds\n", cfg.Load.Enabled, cfg.Load.HighLoadMark, cfg.Load.PeriodSecs)
		ticker := time.NewTicker(time.Duration(cfg.Load.PeriodSecs) * time.Second)
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

// getRunning, setRunning, waitOnRunning and stoppedRunning
// provide synchronization around starting and stopping of the monitor
// there are some tricky edge cases and this ensures only one monitor is running
// for each instance of the loadMonitor and that monitor.Stop() blocks until the monitor
// actually ends
func (lm *loadMonitor) getRunning() bool {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	return lm.isRunning
}

func (lm *loadMonitor) waitOnRunning() {
	lm.wg.Wait()
}

func (lm *loadMonitor) setRunning() {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	if lm.isRunning {
		return
	}
	lm.isRunning = true
	lm.wg.Add(1)
	return
}

func (lm *loadMonitor) stoppedRunning() {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	if !lm.isRunning {
		return
	}
	lm.isRunning = false
	lm.wg.Done()
	return
}

func (lm *loadMonitor) Stop() {
	close(lm.quiter)
	lm.waitOnRunning()
}

func checkLoad(sr statRetriever, highLoadMark float64, onHighLoad func()) (bool, error) {
	loads, err := sr.getLoad()
	high := false
	if err != nil {
		log.Printf("Error checking LoadAvg: %v", err)
		return high, err
	}
	if loads.fiveMinAvg > highLoadMark {
		high = true
		onHighLoad()
	}
	return high, nil
}
