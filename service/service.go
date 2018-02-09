package service

import (
	"github.com/blockassets/cgminer_client"
	"time"
	"os/exec"
	"log"
)

type Commands interface {
	CgmQuit() error
	Reboot()

}
type Command struct {
}

func (*Command) CgmQuit() error {
	clnt := cgminer_client.New("localhost", 4028, 5*time.Second)
	return clnt.Quit()
}

func (*Command)Reboot() {
	time.Sleep(5 * time.Second)
	log.Printf("Reboot Requested")
	exec.Command("/sbin/reboot", "-f").Run()
}
