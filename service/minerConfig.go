package service

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Jeffail/gabs"
)

const configFilePath = "/usr/app/conf.default"

type PoolAddresses struct {
	Pool1 string `json:"pool1"`
	Pool2 string `json:"pool2"`
	Pool3 string `json:"pool3"`
}

func UpdatePools(poolData []byte) error {
	pools := &PoolAddresses{}
	err := json.Unmarshal(poolData, pools)
	if err != nil {
		return err
	}

	jsonParsed, err := gabs.ParseJSONFile(configFilePath)
	if err != nil {
		return err
	}

	bytes := mutateConfig(pools, jsonParsed)

	err = ioutil.WriteFile(configFilePath, bytes,0644)
	if err != nil {
		return err
	}

	return nil
}

func mutateConfig(pools *PoolAddresses, config *gabs.Container) []byte {
	config.Set(pools.Pool1, "pool1")
	config.Set(pools.Pool2, "pool2")
	config.Set(pools.Pool3, "pool3")
	return config.BytesIndent("", "\t")
}
