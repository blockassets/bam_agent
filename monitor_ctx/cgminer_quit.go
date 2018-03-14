package monitor_ctx

import (
	"context"
	"log"

	"github.com/blockassets/bam_agent/tool"
)

type CGMQuitConfig struct {
	Enabled bool                `json:"enabled"`
	Period  tool.RandomDuration `json:"period"`
}

func NewPeriodicCGMQuit(config *CGMQuitConfig, quitFunc func()) Monitor {
	log.Printf("PeriodicCGMQuitMonitor: cgminer quit in: %v", config.Period.Duration)
	monitorFunc := func(ctx context.Context) {
		quitFunc()
	}
	return &Periodic{config.Enabled, config.Period.Duration, monitorFunc}
}