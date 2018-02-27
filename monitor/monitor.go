package monitor

import (
	"log"
)

type Monitor struct {
	Load LoadConfig `json:"load"`
}

func StartMonitors(cfg *Monitor) {
	// Startup the goroutines to do the stuff that needs to be monitored
	sr := LinuxStatRetriever{}

	log.Println("Monitors being started")
	go monitorLoad(&cfg.Load, sr)
}
