package monitor

import (
	"github.com/blockassets/bam_agent/controller"
	"log"
	"time"
)

func monitor_load(sr stat_retrieve, interval time.Duration) {
	for {
		high, err := check_loadAvg(sr)
		if (err == nil) && high {
			controller.Reboot()
		}
		time.Sleep(interval)
	}
}

func check_loadAvg(sr stat_retrieve) (high bool, err error) {
	loads, err := sr.getLoad()
	high = false
	if err == nil {
		//the array loads has three values in it. the 1 min, 5 min and 15 min loadaverage
		if loads[1] > 5.0 {
			high = true
		}
	} else {
		log.Println("Monitor load error: %s", err)
	}
	return high, err
}
