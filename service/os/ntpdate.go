package os

import (
	"log"
	"os/exec"

	"go.uber.org/fx"
)

type Ntpdate interface {
	Ntpdate() error
}

type NtpdateData struct {
	run func(cmd string, arg ...string) error
}

func (r *NtpdateData) Ntpdate() error {
	log.Printf("ntpdate requested")
	return r.run("/usr/bin/ntpdate", "-u", "time.google.com")
}

var NtpdateModule = fx.Provide(func() Ntpdate {
	return &NtpdateData{
		run: func(cmd string, arg ...string) error {
			return exec.Command(cmd, arg...).Run()
		},
	}
})
