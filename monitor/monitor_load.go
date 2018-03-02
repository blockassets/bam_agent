package monitor

import (
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type HighLoadConfig struct {
	Enabled         bool    `json:"enabled"`
	PeriodInSeconds int     `json:"periodInSeconds"`
	HighLoadMark    float64 `json:"highLoadMark"`
}

// Implements the Monitor interface
type LoadMonitor struct {
	*Context
	config        *HighLoadConfig
	statRetriever *service.StatRetriever
	tickerPeriod  *time.Duration
	onHighLoad    func()
}

func newLoadMonitor(context *Context, config *HighLoadConfig, tickerPeriod *time.Duration, statRetriever service.StatRetriever, onHighLoad func()) Monitor {
	return &LoadMonitor{
		Context:       context,
		config:        config,
		tickerPeriod:  tickerPeriod,
		statRetriever: &statRetriever,
		onHighLoad:    onHighLoad,
	}
}

func (monitor *LoadMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("LoadMonitor: Checking load > %v every %v seconds\n", monitor.config.HighLoadMark, monitor.config.PeriodInSeconds)

		go func() {
			monitor.waitGroup.Add(1)
			ticker := time.NewTicker(*monitor.tickerPeriod)
			for {
				select {
				case <-ticker.C:
					checkLoad(*monitor.statRetriever, monitor.config.HighLoadMark, monitor.onHighLoad)
				case <-monitor.quit:
					ticker.Stop()
					monitor.waitGroup.Done()
					return
				}
			}
		}()
	} else {
		log.Println("LoadMonitor: Not enabled")
	}

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
