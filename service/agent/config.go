package agent

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sync"

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
	defaultConfig []byte
	cmdLine       tool.CmdLine
	originalData  *gabs.Container
	loadedData    *FileConfig
	sync.Mutex
}

func (cfg *ConfigData) Original() *gabs.Container {
	return cfg.originalData
}

func (cfg *ConfigData) Loaded() *FileConfig {
	return cfg.loadedData
}

func (cfg *ConfigData) Load() error {
	jsonData := loadJsonFile(cfg.cmdLine.AgentConfigPath)

	var err error
	var merged = cfg.defaultConfig

	// Merge our saved data over the json defaults so
	// that we automatically pick up new configuration defaults over time
	if bytes.Compare(jsonData, cfg.defaultConfig) != 0 {
		merged, err = tool.Merge(jsonData, cfg.defaultConfig)
		if err != nil {
			return err
		}
	}

	// Store a copy of the merged data as 'original' data
	cfg.originalData, err = gabs.ParseJSON(merged)
	if err != nil {
		return err
	}

	// Load our config into FileConfig{}, which does json transforms
	return jsoniter.Unmarshal(merged, &cfg.loadedData)
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

var emptyJson = []byte("{}")

func loadJsonFile(path string) []byte {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil || os.IsNotExist(err) {
		return emptyJson
	}

	stat, err := file.Stat()
	if err != nil {
		return emptyJson
	}

	if stat.Size() == 0 {
		return emptyJson
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return emptyJson
	}

	return data
}

func save(path string, bytes []byte) error {
	return ioutil.WriteFile(path, bytes, 0644)
}

/*
	Returns a loaded config object based on the
	parameters passed in from the cmdLine.
*/
func NewConfig(cmdLine tool.CmdLine, defaultConfig []byte) Config {
	tool.RegisterJsonTypes()

	cfg := ConfigData{cmdLine: cmdLine, defaultConfig: defaultConfig}
	err := cfg.Load()
	// Should never have this issue
	if err != nil {
		log.Panic(err)
	}
	cfg.Save()

	return &cfg
}

var ConfigModule = fx.Options(
	fx.Provide(func(cmdLine tool.CmdLine, confRiceBox tool.ConfRiceBox) Config {
		data, err := (*confRiceBox).Bytes("bam_agent.json")
		if err != nil {
			log.Panic("could not load bam_agent.json from rice")
			return nil
		}
		return NewConfig(cmdLine, data)
	}),

	ConfigMonitorModule,
	ConfigControllerModule,
	ConfigLocationModule,
)
