package monitor

import (
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type AcceptedConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
}

// Implements the Monitor interface
type AcceptedMonitor struct {
	*Context
	config    *AcceptedConfig
	miner     service.Miner
	onStall   func()
	lastShare int64
}

func newAcceptedMonitor(context *Context, config *AcceptedConfig, miner service.Miner, onStall func()) Monitor {
	return &AcceptedMonitor{
		Context: context,
		config:  config,
		miner:   miner,
		onStall: onStall,
	}
}

func (mon *AcceptedMonitor) Start() error {
	if mon.config.Enabled {
		log.Printf("AcceptedMonitor: Checking shares increasing every %v\n", mon.config.Period)
		// reset the tracking attributes
		mon.lastShare = 0
		go mon.makeTickerFunc(func() {
			stalled, err := mon.checkAcceptedShare()
			if err != nil {
				// do nothing, try again next cycle
				// as the miner could be in middle of a restart
			} else {
				if stalled {
					mon.onStall()
				}
			}
		}, mon.config.Period)()
	} else {
		log.Println("AcceptedMonitor: Not enabled")
	}
	return nil
}

func (mon *AcceptedMonitor) checkAcceptedShare() (bool, error) {
	newShare, err := service.GetAccepted(mon.miner)
	if err != nil {
		return false, err
	}
	stalled := false
	// if its less than, then there was a miner restart, so all is good till next check
	// if its greater than, all is good
	if mon.lastShare == newShare {
		stalled = true
	}
	mon.lastShare = newShare
	return stalled, nil
}
