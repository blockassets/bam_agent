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

// create an interface to the actual miner client, so we can mock it out in testing
type CgClient interface {
	Devs() (*[]cgminer_client.Dev, error)
	Quit() error
	Restart() error
	Summary() (*cgminer_client.Summary, error)
	ChipStat() (*[]cgminer_client.ChipStat, error)
}

// Implements the miner.Client interface
type ClientWrapper struct {
	client CgClient
}

func (c *ClientWrapper) Quit() error {
	return c.client.Quit()
}

func (c *ClientWrapper) GetAccepted() (int64, error) {
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

func (c *ClientWrapper) GetTemp() (float64, error) {
	devs, err := c.client.Devs()
	if err != nil {
		return 0, err
	}
	// Temp is same across all boards
	return (*devs)[0].Temperature, nil
}

func NewCgMinerClient(port int64) *cgminer_client.Client {
	return cgminer_client.New(minerHostname, port, minerTimeout)
}

func NewClientWrapper(client CgClient) *ClientWrapper {
	return &ClientWrapper{client: client}
}

var ClientModule = fx.Options(
	fx.Provide(func(port ConfigPort) *cgminer_client.Client {
		return NewCgMinerClient(port.Get())
	}),

	fx.Provide(func(client *cgminer_client.Client) *ClientWrapper {
		return NewClientWrapper(client)
	}),
)
