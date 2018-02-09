package monitor

import (
	"log"
	"time"
)

const (
	restartMinerHours = 23
	rebootHours       = 71
)

func StartMonitors() {
	// Startup the goroutines to do the stuff that needs to be monitored
	sr := LinuxStatRetriever{}

	log.Println("Monitors being started")
	go monitorLoad(sr, time.Minute) // check for 5min average CPU load every minute
	go monitorCyclicMinerRestart(restartMinerHours)
	go monitorCyclicReboot(rebootHours)
}
