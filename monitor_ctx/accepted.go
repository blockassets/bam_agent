package monitor_ctx

import (
	"context"
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type AcceptedConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
}

func NewAcceptedMonitor(config *AcceptedConfig, miner service.Miner, onStall func()) Monitor {
	log.Printf("AcceptedShareMonitor: Checking share %v\n", config.Period)
	lastShare := int64(0)
	monitorFunc := func(ctx context.Context) {
		stalled, err := checkAcceptedShare(&lastShare, miner)
		if err == nil && stalled {
			onStall()
		}
	}
	return &Periodic{config.Enabled, config.Period, monitorFunc}
}

func checkAcceptedShare(lastShare *int64, miner service.Miner) (bool, error) {
	newShare, err := service.GetAccepted(miner)
	if err != nil {
		*lastShare = 0
		return false, err
	}
	stalled := false
	// if its less than, then there was a miner restart, so all is good till next check
	// if its greater than, all is good
	if *lastShare == newShare {
		stalled = true
	}
	*lastShare = newShare
	return stalled, nil
}
