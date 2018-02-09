package monitor


import (
	"github.com/blockassets/bam_agent/controller"
	"log"
	"math/rand"
	"time"
)

func monitorCyclicReboot() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// We dont want to reset miners too close together... if multiple miners are reset or added at same time,
	// this ensures there is a spread on the time to restart the miner
	timeToWait := time.Duration(71)*time.Hour + time.Duration(r1.Intn(3600))*time.Second
	log.Println("Time to wait before Reboot:", timeToWait)
	time.Sleep(timeToWait)
	timeToWait = time.Duration(72) * time.Hour
	controller.Reboot()
	for {
		time.Sleep(timeToWait)
		controller.Reboot()
	}

}