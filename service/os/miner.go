package os

import (
	"log"
	"os/exec"

	"go.uber.org/fx"
)

type Miner interface {
	Start() error
	Stop() error
}

type MinerData struct {
	run func(cmd string) error
}

func (s MinerData) Start() error {
	log.Printf("CGMiner Service Start Requested")
	return s.run("systemctl start cgminer")
}

func (s MinerData) Stop() error {
	log.Printf("CGMiner Service Stop Requested")
	return s.run("systemctl stop cgminer")
}

func NewMiner() Miner {
	return &MinerData{
		run: func(cmd string) error {
			return exec.Command(cmd).Run()
		},
	}
}

var MinerModule = fx.Provide(NewMiner)
