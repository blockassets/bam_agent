package cgminer

import (
	"time"

	"github.com/blockassets/cgminer_client"
	"go.uber.org/fx"
)

var (
	minerHostname = "localhost"
	minerTimeout  = 5 * time.Second
)

// Implements the miner.client interface
type Wrapper struct {
	client cgminer_client.Client
}

func (c *Wrapper) Quit() error {
	return c.client.Quit()
}

func (c *Wrapper) Restart() error {
	return c.client.Restart()
}

func (c *Wrapper) GetAccepted() (int64, error) {
	devs, err := c.client.Devs()
	if err != nil {
		return 0, err
	}
	accepted := int64(0)
	for _, dev := range *devs {
		accepted += dev.Accepted
	}
	return accepted, nil
}

func (c *Wrapper) GetTemp() (float64, error) {
	devs, err := c.client.Devs()
	if err != nil {
		return 0, err
	}
	// Temp is same across all boards
	return (*devs)[0].Temperature, nil
}

func NewCgMinerClient(port int64) cgminer_client.Client {
	return cgminer_client.New(minerHostname, port, minerTimeout)
}

func NewClientWrapper(client cgminer_client.Client) *Wrapper {
	return &Wrapper{client: client}
}

var ClientModule = fx.Options(
	fx.Provide(func(port ConfigPort) cgminer_client.Client {
		return NewCgMinerClient(port.Get())
	}),

	fx.Provide(func(client cgminer_client.Client) *Wrapper {
		return NewClientWrapper(client)
	}),
)
