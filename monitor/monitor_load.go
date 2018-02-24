package monitor

import (
	"log"
	"time"
	"encoding/json"

	"github.com/blockassets/bam_agent/service"
)

const (
	MINIMUM_TIME_BETWEEN_LOAD_TEST = 5 * time.Second // Hard coded so that if there is a config file snafu, the machines/processes wont get into a state
)

type monitorLoadConfig struct {
	*MonitorConfig
	sr statRetriever
}

func (cfg *monitorLoadConfig) ReloadConfig(blob []byte) {
	// Use an anonymous struct to get to inner json
	c := &struct {
		Monitor struct{ System_load monitorLoadConfig }
	}{}
	c.Monitor.System_load.MonitorConfig = &MonitorConfig{}
	err := json.Unmarshal(blob, c)
	cfg.MonitorConfig = c.Monitor.System_load.MonitorConfig
	if err != nil {
		log.Printf("Load Monitor: error getting config %v", err)
		cfg.Enabled = false
		return
	}
}

func (cfg *monitorLoadConfig) Action() {
	high, err := checkLoadAvg(cfg.sr)
	if (err == nil) && high {
		cmd := service.Command{}
		cmd.Reboot()
	}
}

func monitorLoad(msg_ch chan interface{}, sr statRetriever) {
	cfg := &monitorLoadConfig{}
	cfg.sr = sr
	periodicMonitor(msg_ch, MINIMUM_TIME_BETWEEN_LOAD_TEST, cfg)

}

func checkLoadAvg(sr statRetriever) (bool, error) {
	loads, err := sr.getLoad()
	high := false
	if (err == nil) && (loads.fiveMinAvg > 5.0) {
		high = true
	} else {
		log.Println("Monitor load error:", err)
	}
	return high, err
}
