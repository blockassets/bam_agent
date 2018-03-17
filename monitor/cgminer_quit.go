package monitor

import (
	"context"
	"time"

	"github.com/blockassets/bam_agent/service/miner"
)

type CGMQuitConfig struct {
	Enabled bool
	Period  time.Duration
}

func NewCGMQuitMonitor(config CGMQuitConfig, client miner.Client) Result {
	return Result{
		Monitor: &Data{
			Enabled: config.Enabled,
			Period:  config.Period,
			OnTick: func() TickerFunc {
				return func(ctx context.Context) {
					client.Quit()
				}
			},
		},
	}
}
