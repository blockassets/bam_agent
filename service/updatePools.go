package service

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/Jeffail/gabs"
)

type PoolAddresses struct {
	Pool1 string `json:"pool1"`
	Pool2 string `json:"pool2"`
	Pool3 string `json:"pool3"`
}

func UpdatePools(poolsAsJson io.ReadCloser, configFilePath string) error {
	buf, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	out, err := updateCfgJson(poolsAsJson, buf)
	err = ioutil.WriteFile(configFilePath, out, 0644)
	return err
}

func updateCfgJson(poolsAsJson io.ReadCloser, cfgBuf []byte) ([]byte, error) {
	pools := &PoolAddresses{}
	// new values
	err := json.NewDecoder(poolsAsJson).Decode(&pools)
	if err != nil {
		return nil, err
	}
	// We don't own the miner config file, so use the gabs library to
	// easily read and write the file, only being concerned with the fields that we are interested in
	//
	jsonParsed, err := gabs.ParseJSON(cfgBuf)
	if err != nil {
		return nil, err
	}
	jsonParsed.Set(pools.Pool1, "pool1")
	jsonParsed.Set(pools.Pool2, "pool2")
	jsonParsed.Set(pools.Pool3, "pool3")

	return jsonParsed.BytesIndent("", "\t"), nil
}
