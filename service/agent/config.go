package agent

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/Jeffail/gabs"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

const (
	defaultConfigFile = "bam_agent.json"
)

type Config interface {
	Load() error
	Save() error
	Original() *gabs.Container
	Loaded() *FileConfig
}

type ConfigData struct {
	cmdLine      tool.CmdLine
	originalData *gabs.Container
	loadedData   *FileConfig
}

func (cfg *ConfigData) Original() *gabs.Container {
	return cfg.originalData
}

func (cfg *ConfigData) Loaded() *FileConfig {
	return cfg.loadedData
}

func (cfg *ConfigData) Load() error {
	jsonData, err := loadJson(cfg.cmdLine.AgentConfigPath)
	if err != nil {
		return err
	}

	// Need to keep a copy of the original for modifying and writing out
	cfg.originalData, err = gabs.ParseJSON(jsonData)
	if err != nil {
		return err
	}

	// This transforms fields into data we can use in the app
	err = jsoniter.Unmarshal(jsonData, &cfg.loadedData)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigData) Data() Config {
	return cfg
}

func (cfg *ConfigData) Save() error {
	pretty, err := jsoniter.MarshalIndent(cfg.originalData.Data(), "", "  ")
	if err != nil {
		return err
	}

	err = save(cfg.cmdLine.AgentConfigPath, pretty)
	if err != nil {
		return err
	}

	err = cfg.Load()
	if err != nil {
		return err
	}
	return nil
}

/*
	1. Look for the config file in /etc/bam_agent.conf
	2. If the config doesn't exist, load it from the box.
	3. Attempt to write the config to outputFile.
*/
func loadJson(path string) ([]byte, error) {
	var jsonData []byte

	readOnly, errOpen := os.Open(path)
	defer readOnly.Close()

	// Determine how to get the data, either on disk or load default file
	stat, errStat := readOnly.Stat()
	if errStat != nil || os.IsNotExist(errOpen) || stat.Size() == 0 {
		confBox, err := rice.FindBox("../../conf")
		if err != nil {
			return nil, err
		}

		jsonData, err = confBox.Bytes(defaultConfigFile)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(path, jsonData, 0664)
		if err != nil {
			log.Println("Warning: failed to write default bam_agent config file:", err)
		}
	} else {
		jsonData, errOpen = ioutil.ReadAll(readOnly)
		if errOpen != nil {
			return nil, errOpen
		}
	}

	return jsonData, nil
}

func save(path string, bytes []byte) error {
	return ioutil.WriteFile(path, bytes, 0644)
}

/*
	Returns a loaded config object based on the
	parameters passed in from the cmdLine.
*/
func NewConfig(cmdLine tool.CmdLine) Config {
	tool.RegisterJsonTypes()

	cfg := ConfigData{cmdLine: cmdLine}
	cfg.Load()
	cfg.Save()

	return &cfg
}

var ConfigModule = fx.Options(
	fx.Provide(func(cmdLine tool.CmdLine) Config {
		return NewConfig(cmdLine)
	}),

	ConfigMonitorModule,
	ConfigControllerModule,
)
