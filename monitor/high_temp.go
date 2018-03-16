package monitor

import (
	"context"
	"time"

	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
)

type HighTempConfig struct {
	Enabled  bool
	Period   time.Duration
	HighTemp float64
}

func NewHighTempMonitor(config HighTempConfig, client miner.Client, miner os.Miner) Result {
	return Result{
		Monitor: &Data{
			Enabled: config.Enabled,
			Period:  config.Period,
			OnTick: func() TickerFunc {
				return func(ctx context.Context) {
					temp, err := client.GetTemp()
					if err == nil && temp >= config.HighTemp {
						miner.Stop()
					}
				}
			},
		},
	}
}
