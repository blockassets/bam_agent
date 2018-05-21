package os

import (
	"log"
	"os/exec"

	"go.uber.org/fx"
)

const defaultNtpServer = "time.google.com"
const defaultNtpCmd = "/usr/bin/ntpdate"

type Ntpdate interface {
	Ntpdate(server string) error
}

type NtpdateData struct {
	run func(cmd string, arg ...string) error
}

func (r *NtpdateData) Ntpdate(server string) error {
	if len(server) == 0 {
		server = defaultNtpServer
	}
	log.Printf("%s -u %s", defaultNtpCmd, server)
	return r.run(defaultNtpCmd, "-u", server)
}

var NtpdateModule = fx.Provide(func() Ntpdate {
	return &NtpdateData{
		run: func(cmd string, arg ...string) error {
			return exec.Command(cmd, arg...).Run()
		},
	}
})
