package miner

import (
	"io/ioutil"
	"log"

	"github.com/Jeffail/gabs"
	"github.com/blockassets/bam_agent/tool"
	"go.uber.org/fx"
)

type Config interface {
	Load() error
	Save() error
	Data() *gabs.Container
}

type ConfigData struct {
	CmdLine   tool.CmdLine
	Container *gabs.Container
}

/*
	Returns a loaded config object based on the
	parameters passed in from the cmdLine.
*/
func NewConfig(cmdLine tool.CmdLine) Config {
	cfg := ConfigData{CmdLine: cmdLine}
	cfg.Load()
	return &cfg
}

var ConfigModule = fx.Options(
	fx.Provide(
		func(cmdLine tool.CmdLine) Config {
			return NewConfig(cmdLine)
		},
	),

	// Helpers for reading/mutating the config
	PortModule,
	PoolModule,
	NetworkModule,
	FrequencyModule,
)

func (cfg *ConfigData) Data() *gabs.Container {
	return cfg.Container
}

func save(path string, bytes []byte) error {
	return ioutil.WriteFile(path, bytes, 0644)
}

func (cfg *ConfigData) Load() error {
	container, err := gabs.ParseJSONFile(cfg.CmdLine.MinerConfigPath)
	if err != nil {
		container, err = gabs.ParseJSON([]byte(DefaultConfigFile))
		if err != nil {
			log.Fatalln("failed to parse the default miner configuration.")
		}
	}
	cfg.Container = container
	return nil
}

func (cfg *ConfigData) Save() error {
	return save(cfg.CmdLine.MinerConfigPath, cfg.Data().BytesIndent("", "\t"))
}
