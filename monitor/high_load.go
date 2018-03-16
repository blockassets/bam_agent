package monitor

import (
	"context"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

type HighLoadConfig struct {
	Enabled      bool
	Period       time.Duration
	HighLoadMark float64
}

func NewLoadMonitor(config HighLoadConfig, retriever os.StatRetriever, reboot os.Reboot) Result {
	return Result{
		Monitor: &Data{
			Period:  config.Period,
			Enabled: config.Enabled,
			OnTick: func() TickerFunc {
				return func(ctx context.Context) {
					loads, err := retriever.GetLoadData()
					if err == nil && loads.FiveMinAvg > config.HighLoadMark {
						reboot.Reboot()
					}
				}
			},
		},
	}
}
