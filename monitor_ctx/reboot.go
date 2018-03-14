package monitor_ctx

import (
	"context"
	"log"

	"github.com/blockassets/bam_agent/tool"
)

type RebootConfig struct {
	Enabled bool                `json:"enabled"`
	Period  tool.RandomDuration `json:"period"`
}

func NewPeriodicReboot(config *RebootConfig, rebootFunc func()) Monitor {
	log.Printf("PeriodicRebootMonitor(enabled == %v): reboot in %v ", config.Enabled, config.Period.Duration)
	monitorFunc := func(ctx context.Context) {
		rebootFunc()
	}
	return &Periodic{config.Enabled, config.Period.Duration, monitorFunc}
}
