package monitor

import (
	"log"
	"math/rand"
	"time"
	"github.com/blockassets/bam_agent/service"
)

func monitorCyclicMinerRestart(delayInHours int) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	// We dont want to reset miners too close together... if multiple miners are reset or added at same time,
	// this ensures there is a spread on the time to restart the miner
	timeToWait := time.Duration(delayInHours)*time.Hour + time.Duration(r1.Intn(3600))*time.Second
	log.Println("Time to wait before Miner Restart:", timeToWait)
	time.Sleep(timeToWait)
	timeToWait = time.Duration(24) * time.Hour

	cmds := service.Command{}
	cmds.CgmQuit()
	for {
		time.Sleep(timeToWait)
		cmds.CgmQuit()
	}

}
