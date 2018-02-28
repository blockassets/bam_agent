package monitor

import (
	"log"

	"github.com/blockassets/bam_agent/service"
)

type MonitorConfig struct {
	Load LoadConfig `json:"load"`
}

type Monitor interface {
	Start(cfg *MonitorConfig) error
	Stop()
}

func StartMonitors(cfg *MonitorConfig) {
	// Startup the goroutines to do the stuff that needs to be monitored
	sr := LinuxStatRetriever{}

	log.Println("Monitors being started")

	lm := newLoadMonitor(&sr, service.Reboot)
	lm.Start(cfg)
}
