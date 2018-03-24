package miner

import (
	"time"

	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"github.com/blockassets/cgminer_client"
	"go.uber.org/fx"
)

var (
	minerHostname = "localhost"
	minerTimeout  = 5 * time.Second
)

type Client interface {
	Quit() error
	GetAccepted() (int64, error)
	GetTemp() (float64, error)
}

// TODO: In the future, we will need to inject the Port from another source. whichever one isn't null is the one that we use
func NewClient(port cgminer.ConfigPort) Client {
	c := cgminer_client.New(minerHostname, port.Get(), minerTimeout)
	return &cgminer.ClientData{Client: c}
}

var ClientModule = fx.Provide(func(port cgminer.ConfigPort) Client {
	return NewClient(port)
})
