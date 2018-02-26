package service

import (
	"log"
	"os/exec"
	"time"
)

func Reboot() {
	time.Sleep(5 * time.Second)
	log.Printf("Reboot Requested")
	exec.Command("/sbin/reboot", "-f").Run()
}
