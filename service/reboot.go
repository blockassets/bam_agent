package service

import (
	"log"
	"os/exec"
	"time"
)

const (
	DELAY_BEFORE_REBOOT = 5 * time.Second
)

func Reboot() {
	// Give enough time for any calls to return to client before rebooting
	time.Sleep(DELAY_BEFORE_REBOOT)
	log.Printf("Reboot Requested")
	exec.Command("/sbin/reboot", "-f").Run()
}
