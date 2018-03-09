package service

import (
	"log"

	"github.com/blockassets/cgminer_client"
)

func GetAccepted(client *cgminer_client.Client) int64 {
	devs, err := client.Devs()
	if err != nil {
		log.Printf("Error getting accepted shares: %v", err)
		return 0
	}
	accepted := int64(0)
	for _, dev := range *devs {
		accepted += dev.Accepted
	}
	return accepted
}
