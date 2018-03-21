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
	run func(cmd string, arg string) error
	sync func(cmd string) error
}

func (r *RebootData) Reboot() error {
	log.Printf("Reboot Requested")
	r.sync("/bin/sync")
	r.sync("/bin/sync")
	return r.run("/sbin/reboot", "-f")
}

var RebootModule = fx.Provide(func() Reboot {
	return &RebootData{
		run: func(cmd string, arg string) error {
			return exec.Command(cmd, arg).Run()
		},
		sync: func(cmd string) error {
			return exec.Command(cmd).Run()
		},
	}
})
