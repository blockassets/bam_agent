package monitor

import (
	"github.com/blockassets/bam_agent/controller"
	"log"
	"time"
)

func monitorLoad(sr statRetrieve, interval time.Duration) {
	for {
		high, err := checkLoadAvg(sr)
		if (err == nil) && high {
			controller.Reboot()
		}
		time.Sleep(interval)
	}
}

func checkLoadAvg(sr statRetrieve) (bool, error) {
	loads, err := sr.getLoad()
	high := false
	if err == nil {
		if loads.fiveMinAvg > 5.0 {
			high = true
		}
	} else {
		log.Println("Monitor load error: %v", err)
	}
	return high, err
}
