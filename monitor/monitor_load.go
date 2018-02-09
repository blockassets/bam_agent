package monitor

import (
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

func monitorLoad(sr statRetriever, interval time.Duration) {
	for {
		high, err := checkLoadAvg(sr)
		if (err == nil) && high {
			cmds := service.Command{}
			cmds.Reboot()
		}
		time.Sleep(interval)
	}
}

func checkLoadAvg(sr statRetriever) (bool, error) {
	loads, err := sr.getLoad()
	high := false
	if (err == nil) && (loads.fiveMinAvg > 5.0) {
		high = true
	} else {
		log.Printf("Monitor load error: %v", err)
	}
	return high, err
}
