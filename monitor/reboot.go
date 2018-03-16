package monitor

import (
	"context"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

type RebootConfig struct {
	Enabled bool
	Period  time.Duration
}

func NewRebootMonitor(config RebootConfig, reboot os.Reboot) Result {
	return Result{
		Monitor: &Data{
			Enabled: config.Enabled,
			Period:  config.Period,
			OnTick: func() TickerFunc {
				return func(ctx context.Context) {
					reboot.Reboot()
				}
			},
		},
	}
}
