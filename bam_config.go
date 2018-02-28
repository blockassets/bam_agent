package main

import (
	"io/ioutil"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/blockassets/bam_agent/monitor"
	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigDefault
)

type BamConfig struct {
	Monitor monitor.MonitorConfig `json:"monitor"`
}

func InitialiseConfigFile(configFile string) (*BamConfig, error) {

	confBox, err := rice.FindBox("conf")
	if err != nil {
		return nil, err
	}
	defaultJson, err := confBox.Bytes("config.json")
	// create a config structure from the default json
	// so any struct additions in this version of app will get correct
	// defaults
	bamConfig := &BamConfig{}
	err = json.Unmarshal(defaultJson, bamConfig)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0664)
	if err == nil {
		// write out default content if just created
		file.Write(defaultJson)
		file.Close()
	} else {
		if os.IsExist(err) {
			file, err = os.Open(configFile)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			buf, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(buf, bamConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	return bamConfig, nil
}
