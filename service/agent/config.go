package agent

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"

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
	Update(path string, data interface{}) error
}

type ConfigData struct {
	cmdLine      tool.CmdLine
	originalData *gabs.Container
	loadedData   *FileConfig
	sync.Mutex
}

func (cfg *ConfigData) Original() *gabs.Container {
	return cfg.originalData
}

func (cfg *ConfigData) Loaded() *FileConfig {
	return cfg.loadedData
}

func (cfg *ConfigData) Load() error {
	// Load both the default file and the config file
	jsonData, errFile := loadJsonFile(cfg.cmdLine.AgentConfigPath)
	jsonDefaults, errDefaults := loadJsonDefaults()

	// Should never reach this since defaults load from the box
	if errDefaults != nil {
		return errDefaults
	}

	// No config file, so use the defaults
	if errFile != nil {
		jsonData = jsonDefaults
	}

	// Merge our saved data over the json defaults so
	// that we automatically pick up new configuration defaults over time
	mergedStr, err := tool.Merge(jsonData, jsonDefaults)
	if err != nil {
		return err
	}

	// Store a copy of the merged data as 'original' data
	cfg.originalData, err = gabs.ParseJSON(mergedStr)
	if err != nil {
		return err
	}

	// Load our config into FileConfig{}, which does json transforms
	return jsoniter.Unmarshal(mergedStr, &cfg.loadedData)
}

func (cfg *ConfigData) Data() Config {
	return cfg
}

func (cfg *ConfigData) Update(path string, data interface{}) error {
	cfg.Lock()
	defer cfg.Unlock()

	err := configUpdate(cfg.Original(), path, data)
	if err != nil {
		return err
	}
	return cfg.Save()
}

func configUpdate(original *gabs.Container, path string, data interface{}) error {
	converted, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}

	updated, err := gabs.ParseJSON(converted)
	if err != nil {
		return err
	}

	_, err = original.Set(updated.Data(), path)
	return err
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

func loadJsonFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil || os.IsNotExist(err) {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() == 0 {
		return nil, errors.New("empty file")
	}

	return ioutil.ReadAll(file)
}

func loadJsonDefaults() ([]byte, error) {
	confBox, err := rice.FindBox("../../conf")
	if err != nil {
		return nil, err
	}

	return confBox.Bytes(defaultConfigFile)
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
	err := cfg.Load()
	// Should never have this issue
	if err != nil {
		log.Panic(err)
	}
	cfg.Save()

	return &cfg
}

var ConfigModule = fx.Options(
	fx.Provide(func(cmdLine tool.CmdLine) Config {
		return NewConfig(cmdLine)
	}),

	ConfigMonitorModule,
	ConfigControllerModule,
	ConfigLocationModule,
)
