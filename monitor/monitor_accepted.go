package monitor

import (
	"log"
	"time"
)

type AcceptedConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
}

// Implements the Monitor interface
type AcceptedMonitor struct {
	*Context
	config      *AcceptedConfig
	getAccepted func() int64
	onStall     func()
}

func newAcceptedMonitor(context *Context, config *AcceptedConfig, getAccepted func() int64, onStall func()) Monitor {
	return &AcceptedMonitor{
		Context:     context,
		config:      config,
		getAccepted: getAccepted,
		onStall:     onStall,
	}
}

func (monitor *AcceptedMonitor) Start() error {
	if monitor.config.Enabled {
		log.Printf("AcceptedMonitor: Checking shares increasing every %v\n", monitor.config.Period)
		lastShare := int64(0)
		go monitor.makeTickerFunc(func() {
			newShare := monitor.getAccepted()
			if lastShare == newShare {
				monitor.onStall()
			}
			lastShare = newShare
		}, monitor.config.Period)()
	} else {
		log.Println("AcceptedMonitor: Not enabled")
	}

	return nil
}
