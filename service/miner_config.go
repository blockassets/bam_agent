package service

import (
	"io/ioutil"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

const configFilePath = "/usr/app/conf.default"

type PoolAddresses struct {
	Pool1 string `json:"pool1"`
	Pool2 string `json:"pool2"`
	Pool3 string `json:"pool3"`
}

type IPAddresses struct {
	IPAddress string `json:"address"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	Dns       string `json:"dns"`
}

// Local private cache. Always reference this through the LoadMinerConfig() function
var config *gabs.Container

func LoadMinerConfig() (*gabs.Container, error) {
	// Local cache of config to prevent a lot of reads
	if config != nil {
		return config, nil
	}

	config, err := gabs.ParseJSONFile(configFilePath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func SaveMinerConfig(bytes []byte) error {
	err := ioutil.WriteFile(configFilePath, bytes, 0644)
	// Clear the cache so that the next LoadMinerConfig() will read from disk
	if err == nil {
		config = nil
	}
	return err
}

func UpdatePools(poolData []byte) error {
	pools := &PoolAddresses{}
	err := jsoniter.Unmarshal(poolData, pools)
	if err != nil {
		return err
	}

	config, err := LoadMinerConfig()
	if err != nil {
		return err
	}

	bytes := mutatePools(pools, config)
	return SaveMinerConfig(bytes)
}

func mutatePools(pools *PoolAddresses, config *gabs.Container) []byte {
	config.Set(pools.Pool1, "pool1")
	config.Set(pools.Pool2, "pool2")
	config.Set(pools.Pool3, "pool3")
	return config.BytesIndent("", "\t")
}

func GetPools() (*PoolAddresses, error) {
	config, err := LoadMinerConfig()
	if err != nil {
		return nil, err
	}

	pool1, _ := config.Path("pool1").Data().(string)
	pool2, _ := config.Path("pool2").Data().(string)
	pool3, _ := config.Path("pool3").Data().(string)

	return &PoolAddresses{
		Pool1: pool1,
		Pool2: pool2,
		Pool3: pool3,
	}, nil
}

func UpdateIPAddresses(ipData []byte) error {
	ipa := &IPAddresses{}
	err := jsoniter.Unmarshal(ipData, ipa)
	if err != nil {
		return err
	}

	config, err := LoadMinerConfig()
	if err != nil {
		return err
	}

	// Set the /etc/network/interfaces
	err = SetStaticIP(ipa.IPAddress, ipa.Netmask, ipa.Gateway)
	if err != nil {
		return err
	}
	// set the miner config
	bytes := mutateIPAddresses(ipa, config)
	return SaveMinerConfig(bytes)
}

func mutateIPAddresses(ipa *IPAddresses, config *gabs.Container) []byte {
	config.Set(ipa.IPAddress, "ip")
	config.Set(ipa.Netmask, "mask")
	config.Set(ipa.Gateway, "gateway")
	config.Set(ipa.Dns, "dns")

	return config.BytesIndent("", "\t")
}
