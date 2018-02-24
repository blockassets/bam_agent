package service

import (
	"log"
	"time"

	"github.com/blockassets/cgminer_client"
)

const (
	MINER_HOSTNAME = "localhost"
	MINER_PORT     = 4028
	MINER_TIMEOUT  = 5 * time.Second
)

func (*Command) CgmQuit() error {
	log.Printf("cgminer quit requested")
	clnt := cgminer_client.New(MINER_HOSTNAME, MINER_PORT, MINER_TIMEOUT)
	return clnt.Quit()
}
