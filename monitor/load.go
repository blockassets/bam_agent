package monitor

import (
	"context"
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type HighLoadConfig struct {
	Enabled      bool          `json:"enabled"`
	Period       time.Duration `json:"period"`
	HighLoadMark float64       `json:"highLoadMark"`
}

func NewLoadMonitor(config *HighLoadConfig, sr service.StatRetriever, onHighLoad func()) Monitor {
	log.Printf("LoadMonitor(enabled == %v): Checking load > %v every %v\n", config.Enabled, config.HighLoadMark, config.Period)

	monitorFunc := func(ctx context.Context) { checkLoad(sr, config.HighLoadMark, onHighLoad) }

	return &Periodic{config.Enabled, config.Period, monitorFunc}
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
