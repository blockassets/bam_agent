package os

import (
	"log"
	"os/exec"

	"go.uber.org/fx"
)

type Reboot interface {
	Reboot() error
}

type RebootData struct {
	run  func(cmd string, arg string) error
}

func (r *RebootData) Reboot() error {
	log.Printf("Reboot Requested")
	r.run("/bin/sync", "")
	r.run("/bin/sync", "")
	return r.run("/bin/systemctl", "reboot")
}

var RebootModule = fx.Provide(func() Reboot {
	return &RebootData{
		run: func(cmd string, arg string) error {
			return exec.Command(cmd, arg).Run()
		},
	}
})
