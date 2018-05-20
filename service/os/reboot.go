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
}

func (r *RebootData) Reboot() error {
	log.Printf("Reboot Requested")
	err := r.run("/bin/sync", "")
	if err != nil {
		log.Println(err)
	}
	err = r.run("/bin/sync", "")
	if err != nil {
		log.Println(err)
	}
	return r.run("/sbin/reboot", "-f")
}

var RebootModule = fx.Provide(func() Reboot {
	return &RebootData{
		run: func(cmd string, arg string) error {
			return exec.Command(cmd, arg).Run()
		},
	}
})
