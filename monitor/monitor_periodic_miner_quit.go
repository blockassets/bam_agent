package monitor

import (
	"encoding/json"
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

const (
	MINIMUM_TIME_BETWEEN_QUIT = 60 * time.Second // Hard coded so that if there is a config file snafu, the machines/processes wont get into a state
)

// override the PeriodicAction interface implemented in monitor and defined in PeriodicMonitor

type minerQuitConfig struct {
	*MonitorConfig
}

func (cfg *minerQuitConfig) ReloadConfig(blob []byte) {
	// Rather ugly way of getting to the json
	c := &struct {
		Monitor struct{ Miner_quit minerQuitConfig }
	}{}
	c.Monitor.Miner_quit.MonitorConfig = &MonitorConfig{}
	err := json.Unmarshal(blob, c)
	*cfg = c.Monitor.Miner_quit
	if err != nil {
		log.Printf("MinerQuit Monitor: error getting config %v", err)
		cfg.Enabled = false
		return
	}
}

func (*minerQuitConfig) Action() {
	cmds := service.Command{}
	cmds.CgmQuit()
}

func monitorPeriodicMinerQuit(msg_ch chan interface{}) {
	cfg := &minerQuitConfig{}
	periodicMonitor(msg_ch, MINIMUM_TIME_BETWEEN_QUIT, cfg)
}
