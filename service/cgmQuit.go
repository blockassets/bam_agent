package service

import (
	"log"

	"github.com/blockassets/cgminer_client"
)

func CgmQuit(client *cgminer_client.Client) error {
	log.Printf("cgminer quit requested")
	return client.Quit()
}
