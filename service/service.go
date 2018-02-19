package service

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"github.com/blockassets/bam_agent/util"
	"github.com/blockassets/cgminer_client"
)

//TODO: refactor this file to separate out the service level commands into their own files
type Commands interface {
	CgmQuit() error
	Reboot()
	UpdatePools(poolsAsJson io.ReadCloser, configFilePath string) error
}

type Command struct {
}

func (*Command) CgmQuit() error {
	clnt := cgminer_client.New("localhost", 4028, 5*time.Second)
	return clnt.Quit()
}

func (*Command) Reboot() {
	time.Sleep(5 * time.Second)
	log.Printf("Reboot Requested")
	exec.Command("/sbin/reboot", "-f").Run()
}

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
	err = util.UnmarshalJsonObjAndMap(buf, pools, unknown)
	if err != nil {
		return err
	}
	// new values
	err = json.NewDecoder(poolsAsJson).Decode(&pools)
	if err != nil {
		return err
	}
	buf, err = util.MarshalJsonObjAndMap(pools, unknown)
	if err != nil {
		return err
	}
	out := bytes.Buffer{}
	json.Indent(&out, buf, "", "\t")
	err = ioutil.WriteFile(configFilePath, out.Bytes(), 0644)
	return err

}
