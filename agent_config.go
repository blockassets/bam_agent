package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/blockassets/bam_agent/monitor"
	"github.com/json-iterator/go"
)

const (
	configFilePath = "/etc/"
	configFileName = "bam_agent.json"
)

type AgentConfig struct {
	Monitor monitor.Config `json:"monitor"`
}

/*
	1. Look for the config file in /etc/bam_agent.conf
	2. If the config doesn't exist, load it from the box.
	3. Attempt to write the config to outputFile.
	4. Return the parsed json structure.
*/
func LoadAgentConfig(outputFile string) (*AgentConfig, error) {
	var jsonData []byte

	readOnly, err := os.Open(configFilePath + configFileName)
	defer readOnly.Close()

	if os.IsNotExist(err) {
		confBox, err := rice.FindBox("conf")
		if err != nil {
			return nil, err
		}

		jsonData, err = confBox.Bytes(configFileName)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(outputFile, jsonData, 0664)
		if err != nil {
			log.Println("Warning: failed to write default bam_agent config file:", err)
		}
	} else {
		jsonData, err = ioutil.ReadAll(readOnly)
		if err != nil {
			return nil, err
		}
	}

	config := &AgentConfig{}
	err = jsoniter.Unmarshal(jsonData, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
