package miner

import (
	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"go.uber.org/fx"
)

type Client interface {
	Quit() error
	Restart() error
	GetAccepted() (int64, error)
	GetTemp() (float64, error)
}

var ClientModule = fx.Options(
	cgminer.ClientModule,

	// TODO: In the future, we may return different types of clients.
	fx.Provide(func(client *cgminer.Wrapper) Client {
		return client
	}),
)
