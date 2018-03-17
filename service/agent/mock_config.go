package agent

import (
	"log"

	"github.com/GeertJohan/go.rice"
	"github.com/Jeffail/gabs"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

// Type insurance
var _ Config = &MockConfig{}

type MockConfig struct {
	cmdLine   tool.CmdLine
	originalData  *gabs.Container
	loadedData  *FileConfig
	CalledSave bool
}

func (cfg *MockConfig) Original() *gabs.Container {
	return cfg.originalData
}

func (cfg *MockConfig) Loaded() *FileConfig {
	return cfg.loadedData
}

func (cfg *MockConfig) Load() error {
	confBox, err := rice.FindBox("../../conf")
	if err != nil {
		return err
	}

	data, err := confBox.Bytes("bam_agent.json")
	if err != nil {
		return err
	}

	cfg.originalData, err = gabs.ParseJSON(data)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(data, &cfg.loadedData)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *MockConfig) Save() error {
	cfg.CalledSave = true
	err := jsoniter.UnmarshalFromString(cfg.originalData.String(), cfg.loadedData)
	if err != nil {
		log.Panic(err)
	}
	return nil
}

func NewMockConfig() *MockConfig {
	tool.RegisterJsonTypes()
	mc := &MockConfig{}
	err := mc.Load()
	if err != nil {
		log.Panicf("couldn't make mock data: %v", err)
	}
	return mc
}
