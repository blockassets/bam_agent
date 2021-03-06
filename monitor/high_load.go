package monitor

import (
	"context"
	"log"
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
					if err == nil {
						log.Printf("1m load: %v", loads.OneMinAvg)
						if loads.OneMinAvg > config.HighLoadMark {
							err = reboot.Reboot()
							if err != nil {
								log.Println(err)
							}
						}
					} else {
						log.Println(err)
					}
				}
			},
		},
	}
}
