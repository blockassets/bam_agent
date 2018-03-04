package monitor

import (
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type HighLoadConfig struct {
	Enabled      bool          `json:"enabled"`
	Period       time.Duration `json:"period"`
	HighLoadMark float64       `json:"highLoadMark"`
}

// Implements the Monitor interface
type LoadMonitor struct {
	*Context
	config        *HighLoadConfig
	statRetriever *service.StatRetriever
	onHighLoad    func()
}

func newLoadMonitor(context *Context, config *HighLoadConfig, statRetriever service.StatRetriever, onHighLoad func()) Monitor {
	return &LoadMonitor{
		Context:       context,
		config:        config,
		statRetriever: &statRetriever,
		onHighLoad:    onHighLoad,
	}
}

func (monitor *LoadMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("LoadMonitor: Checking load > %v every %v\n", monitor.config.HighLoadMark, monitor.config.Period)

		go monitor.makeTickerFunc(func() {
			checkLoad(*monitor.statRetriever, monitor.config.HighLoadMark, monitor.onHighLoad)
		}, monitor.config.Period)()
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
