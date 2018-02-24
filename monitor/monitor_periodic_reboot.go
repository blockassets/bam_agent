package monitor

import (
	"encoding/json"
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

const (
	MINIMUM_TIME_BETWEEN_REBOOTS = 1 * time.Hour // Hard coded so that if there is a config file snafu, the machines/processes wont get into a state
)

type systemRebootConfig struct {
	*MonitorConfig
}

func (cfg *systemRebootConfig) ReloadConfig(blob []byte) {
	// Use an anonymous struct to get to inner json
	c := &struct {
		Monitor struct{ System_reboot systemRebootConfig }
	}{}
	c.Monitor.System_reboot.MonitorConfig = &MonitorConfig{}
	err := json.Unmarshal(blob, c)
	*cfg = c.Monitor.System_reboot
	if err != nil {
		log.Printf("Reboot Monitor: error getting config %v", err)
		cfg.Enabled = false
		return
	}
}

func (*systemRebootConfig) Action() {
	cmds := service.Command{}
	cmds.Reboot()
}

func monitorPeriodicReboot(msg_ch chan interface{}) {
	cfg := &systemRebootConfig{}
	periodicMonitor(msg_ch, MINIMUM_TIME_BETWEEN_REBOOTS, cfg)
}
