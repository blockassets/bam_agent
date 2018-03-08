package service

import (
	"io/ioutil"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

const configFilePath = "/usr/app/conf.default"

type PoolAddresses struct {
	Pool1 string `json:"pool1"`
	Pool2 string `json:"pool2"`
	Pool3 string `json:"pool3"`
}

type StaticNetConfig struct {
	IPAddress string `json:"address"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	Dns       string `json:"dns"`
}

type MinerFrequency struct {
	Frequency int `json:"frequency"`
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

func UpdateStaticNetConfig(ipData []byte) error {
	netConfig := &StaticNetConfig{}
	err := jsoniter.Unmarshal(ipData, netConfig)
	if err != nil {
		return err
	}

	config, err := LoadMinerConfig()
	if err != nil {
		return err
	}

	// Set the /etc/network/interfaces
	err = SetInterfaceToStaticIP(netConfig.IPAddress, netConfig.Netmask, netConfig.Gateway)
	if err != nil {
		return err
	}
	// set the miner config
	bytes := mutateStaticNetConfig(netConfig, config)
	return SaveMinerConfig(bytes)
}

func mutateStaticNetConfig(netConfig *StaticNetConfig, config *gabs.Container) []byte {
	config.Set(false, "autoNet")
	config.Set(netConfig.IPAddress, "ip")
	config.Set(netConfig.Netmask, "mask")
	config.Set(netConfig.Gateway, "gateway")
	config.Set(netConfig.Dns, "dns")

	return config.BytesIndent("", "\t")
}

func UpdateDHCPNetConfig() error {
	config, err := LoadMinerConfig()
	if err != nil {
		return err
	}

	// Set the /etc/network/interfaces
	err = SetInterfaceToDhcp()
	if err != nil {
		return err
	}
	// set the miner config
	bytes := mutateDHCPNetConfig(config)
	return SaveMinerConfig(bytes)
}

func mutateDHCPNetConfig(config *gabs.Container) []byte {
	config.Set(true, "autoNet")
	config.Set("", "ip")
	config.Set("", "mask")
	config.Set("", "gateway")
	config.Set("", "dns")

	return config.BytesIndent("", "\t")
}

func UpdateFrequency(mfData []byte) error {
	mf := &MinerFrequency{}
	err := jsoniter.Unmarshal(mfData, mf)
	if err != nil {
		return err
	}

	config, err := LoadMinerConfig()
	if err != nil {
		return err
	}

	// set the miner config
	bytes := mutateFrequencyConfig(mf, config)
	return SaveMinerConfig(bytes)
}

func mutateFrequencyConfig(mf *MinerFrequency, config *gabs.Container) []byte {
	config.Set(strconv.Itoa(mf.Frequency), "frequency")
	return config.BytesIndent("", "\t")
}
