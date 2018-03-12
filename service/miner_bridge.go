package service

import (
	"log"
	"os/exec"

	"github.com/blockassets/cgminer_client"
)

//
//	Interface to decouple cgminer_client implementation. Helps in testing.
//
type Miner interface {
	Devs() (*[]cgminer_client.Dev, error)
	Quit() error
}

// Helper functions to derive facts from the miner
// They are unit tested as part of monitor_accepted_test and monitor_high_temp_test
func GetAccepted(miner Miner) (int64, error) {
	devs, err := miner.Devs()
	if err != nil {
		log.Printf("Error getting accepted shares: %v", err)
		return 0, err
	}
	accepted := int64(0)
	for _, dev := range *devs {
		accepted += dev.Accepted
	}
	return accepted, nil
}

func GetTemp(miner Miner) (float64, error) {
	devs, err := miner.Devs()
	if err != nil {
		log.Printf("Error getting Temp: %v", err)
		return 0, err
	}
	// Temp is same across all boards
	// so grab first
	return (*devs)[0].Temperature, nil
}

// Distinct from quiting... this stops it at the Linux service level so it wont
// automatically restart
func StopMiner() {
	log.Printf("CGMiner Service Stop Requested")
	exec.Command("systemctl stop cgminer").Run()
}
