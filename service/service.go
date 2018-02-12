package service

import (
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/blockassets/cgminer_client"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
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
	var pools PoolAddresses
	err := json.NewDecoder(poolsAsJson).Decode(&pools)
	if err == nil {
		// use the gabs json helper..
		cont := gabs.New()
		cont, err = gabs.ParseJSONFile(configFilePath)
		if err == nil {
			cont.SetP(pools.Pool1, "pool1")
			cont.SetP(pools.Pool2, "pool2")
			cont.SetP(pools.Pool3, "pool3")
			err = ioutil.WriteFile(configFilePath, cont.Bytes(), 0644)
		}
	}
	if err != nil {
		log.Println(err)
	}
	return err
}
