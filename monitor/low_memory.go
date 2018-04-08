package monitor

import (
	"context"
	"log"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

type LowMemoryConfig struct {
	Enabled   bool
	Period    time.Duration
	LowMemory float64
}

func NewLowMemoryMonitor(config LowMemoryConfig, memInfo os.MemInfo, reboot os.Reboot) Result {
	return Result{
		Monitor: &Data{
			Period:  config.Period,
			Enabled: config.Enabled,
			OnTick: func() TickerFunc {
				return func(ctx context.Context) {
					loads, err := memInfo.Get()
					if err == nil {
						if avail, ok := loads[os.MemAvailable]; ok {
							log.Printf("MemAvilable: %v", avail)
							if avail < config.LowMemory {
								err = reboot.Reboot()
								if err != nil {
									log.Println(err)
								}
							}
						} else {
							log.Println("Can't find MemAvailable in data")
						}
					} else {
						log.Println(err)
					}
				}
			},
		},
	}
}
