package service

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/blockassets/bam_agent/util"
)

type PoolAddresses struct {
	Pool1 string `json:"pool1"`
	Pool2 string `json:"pool2"`
	Pool3 string `json:"pool3"`
}

func (*Command) UpdatePools(poolsAsJson io.ReadCloser, configFilePath string) error {
	pools := &PoolAddresses{}
	unknown := map[string]json.RawMessage{}
	buf := []byte{}
	buf, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	err = util.UnmarshalJson(buf, pools, unknown)
	if err != nil {
		return err
	}
	// new values
	err = json.NewDecoder(poolsAsJson).Decode(&pools)
	if err != nil {
		return err
	}
	buf, err = util.MarshalJson(pools, unknown)
	if err != nil {
		return err
	}
	out := bytes.Buffer{}
	json.Indent(&out, buf, "", "\t")
	err = ioutil.WriteFile(configFilePath, out.Bytes(), 0644)
	return err
}
