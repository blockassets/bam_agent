package monitor_ctx

import (
	"context"
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type HighTempConfig struct {
	Enabled  bool          `json:"enabled"`
	Period   time.Duration `json:"period"`
	HighTemp float64       `json:"highTemp"`
}

func NewHighTempMonitor(config *HighTempConfig, miner service.Miner, onHighTemp func()) Monitor {
	log.Printf("HighTempMonitor: Checking for temp over %v every %v\n", config.HighTemp, config.Period)
	monitorFunc := func(ctx context.Context) {
		overTemp, err := checkHighTemp(miner, config.HighTemp)
		if err == nil && overTemp {
			onHighTemp()
		}
	}
	return &Periodic{config.Enabled, config.Period, monitorFunc}
}

func checkHighTemp( miner service.Miner, highTemp float64) (bool, error) {
	temp, err := service.GetTemp(miner)
	if err != nil {
		return false, err
	}
	if temp < highTemp {
		return false, nil
	}
	return true, nil
}
