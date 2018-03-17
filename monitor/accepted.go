package monitor

import (
	"context"
	"time"

	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
)

type AcceptedConfig struct {
	Enabled bool
	Period  time.Duration
}

func NewAcceptedMonitor(config AcceptedConfig, client miner.Client, reboot os.Reboot) Result {
	return Result{
		Monitor: &Data{
			Period:  config.Period,
			Enabled: config.Enabled,
			OnTick: func() TickerFunc {
				lastAccepted := int64(0)
				return func(ctx context.Context) {
					currentAccepted, err := client.GetAccepted()
					if err != nil {
						lastAccepted = 0
					} else if lastAccepted > 0 && currentAccepted == lastAccepted {
						reboot.Reboot()
					} else {
						lastAccepted = currentAccepted
					}
				}
			},
		},
	}
}
