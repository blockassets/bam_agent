package monitor

import (
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type LoadConfig struct {
	Enabled      bool    `json:"enabled"`
	PeriodSecs   int     `json:"period_secs"`
	HighLoadMark float64 `json:"high_load_mark"`
}

func monitorLoad(cfg *LoadConfig, sr statRetriever) {
	if cfg.Enabled {
		for {
			high, err := checkLoadAvg(sr, cfg.HighLoadMark)
			if (err == nil) && high {
				service.Reboot()
			}
			time.Sleep(time.Duration(cfg.PeriodSecs) * time.Second)
		}
	}
}

func checkLoadAvg(sr statRetriever, highLoadMark float64) (bool, error) {
	loads, err := sr.getLoad()
	high := false
	if (err == nil) && (loads.fiveMinAvg > highLoadMark) {
		high = true
	} else {
		log.Printf("Monitor load error: %v", err)
	}
	return high, err
}
