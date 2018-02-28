package service

import (
	"log"
	"os/exec"
	"time"
)



func Reboot() {
	log.Printf("Reboot Requested")
	exec.Command("/sbin/reboot", "-f").Run()
}
