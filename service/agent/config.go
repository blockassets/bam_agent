package agent

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/GeertJohan/go.rice"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

type Config struct {
	CmdLine    tool.CmdLine     `json:"cmdLine"`
	Monitor    MonitorConfig    `json:"monitor"`
	Controller ControllerConfig `json:"controller"`
}

/*
	1. Look for the config file in /etc/bam_agent.conf
	2. If the config doesn't exist, load it from the box.
	3. Attempt to write the config to outputFile.
	4. Return the parsed json structure.
*/
func (cfg *Config) Load() error {
	var jsonData []byte

	readOnly, err := os.Open(cfg.CmdLine.AgentConfigPath)
	defer readOnly.Close()

	if os.IsNotExist(err) {
		confBox, err := rice.FindBox("../../conf")
		if err != nil {
			return err
		}

		_, file := filepath.Split(cfg.CmdLine.AgentConfigPath)
		jsonData, err = confBox.Bytes(file)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cfg.CmdLine.AgentConfigPath, jsonData, 0664)
		if err != nil {
			log.Println("Warning: failed to write default bam_agent config file:", err)
		}
	} else {
		jsonData, err = ioutil.ReadAll(readOnly)
		if err != nil {
			return err
		}
	}

	err = jsoniter.Unmarshal(jsonData, &cfg)
	if err != nil {
		return err
	}

	return nil
}

/*
	Returns a loaded config object based on the
	parameters passed in from the cmdLine.
*/
func NewConfig(cmdLine tool.CmdLine) (Config, error) {
	tool.RegisterJsonTypes()

	cfg := &Config{CmdLine: cmdLine}
	err := cfg.Load()

	return *cfg, err
}

var ConfigModule = fx.Options(
	fx.Provide(func(cmdLine tool.CmdLine) (Config, error) {
		return NewConfig(cmdLine)
	}),
)
