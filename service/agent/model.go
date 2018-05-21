package agent

import (
	"time"

	"github.com/blockassets/bam_agent/tool"
)

// Top level config for the bam_agent.json file
type FileConfig struct {
	CmdLine    tool.CmdLine     `json:"cmdLine"`
	Location   LocationConfig   `json:"location"`
	Monitor    MonitorConfig    `json:"monitor"`
	Controller ControllerConfig `json:"controller"`
}

// Location
type LocationConfig struct {
	Facility string `json:"facility"`
	Rack     string `json:"rack"`
	Row      string `json:"row"`
	Shelf    int    `json:"shelf"`
	Position int    `json:"position"`
}

// Controller
type ControllerConfig struct {
	Reboot ControllerRebootConfig `json:"reboot"`
}

type ControllerRebootConfig struct {
	Delay time.Duration `json:"delay"`
}

// Monitor
type MonitorConfig struct {
	HighLoad       HighLoadConfig  `json:"highLoad"`
	HighTemp       HighTempConfig  `json:"highTemperature"`
	AcceptedShares AcceptedConfig  `json:"acceptedShares"`
	CGMQuit        CGMQuitConfig   `json:"cgMinerQuit"`
	Reboot         RebootConfig    `json:"reboot"`
	LowMemory      LowMemoryConfig `json:"lowMemory"`
	Ntpdate        NtpdateConfig   `json:"ntpdate"`
}

type RebootConfig struct {
	Enabled bool                `json:"enabled"`
	Period  tool.RandomDuration `json:"period"`
}

type NtpdateConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
	Server  string        `json:"server"`
}

type HighLoadConfig struct {
	Enabled      bool          `json:"enabled"`
	Period       time.Duration `json:"period"`
	HighLoadMark float64       `json:"highLoadMark"`
}

type AcceptedConfig struct {
	Enabled bool          `json:"enabled"`
	Period  time.Duration `json:"period"`
}

type HighTempConfig struct {
	Enabled  bool          `json:"enabled"`
	Period   time.Duration `json:"period"`
	HighTemp float64       `json:"highTemp"`
}

type CGMQuitConfig struct {
	Enabled bool                `json:"enabled"`
	Period  tool.RandomDuration `json:"period"`
}

type LowMemoryConfig struct {
	Enabled   bool          `json:"enabled"`
	Period    time.Duration `json:"period"`
	LowMemory float64       `json:"lowMemory"`
}
