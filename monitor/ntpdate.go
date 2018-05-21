package monitor

import (
	"context"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

type NtpdateConfig struct {
	Enabled bool
	Period  time.Duration
	Server  string
}

func NewNtpdateMonitor(config NtpdateConfig, ntpdate os.Ntpdate) Result {
	return Result{
		Monitor: &Data{
			Enabled: config.Enabled,
			Period:  config.Period,
			OnTick: func() TickerFunc {
				return func(ctx context.Context) {
					ntpdate.Ntpdate(config.Server)
				}
			},
		},
	}
}
