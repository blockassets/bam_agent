package miner

import (
	"time"

	"github.com/blockassets/cgminer_client"
	"go.uber.org/fx"
)

var (
	minerHostname = "localhost"
	minerTimeout  = 5 * time.Second
)

type Client interface {
	Quit() error
	Devs() (*[]Dev, error)
	GetAccepted() (int64, error)
	GetTemp() (float64, error)
}

type Dev struct {
	Accepted    int64   `json:"Accepted"`
	Temperature float64 `json:"Temperature"`
}

//type BWMinerClient struct {
//	client *cgminer_client.Client
//}

type CGMinerClient struct {
	client *cgminer_client.Client
}

func (c CGMinerClient) Quit() error {
	return c.client.Quit()
}

func (c CGMinerClient) Devs() (*[]Dev, error) {
	clientDevs, err := c.Devs()
	if err != nil {
		return nil, err
	}

	devs := make([]Dev, len(*clientDevs))
	for idx, d := range *clientDevs {
		// TODO: expose more fields
		devs[idx].Accepted = d.Accepted
		devs[idx].Temperature = d.Temperature
	}

	return &devs, nil
}

func (c CGMinerClient) GetAccepted() (int64, error) {
	devs, err := c.Devs()
	if err != nil {
		return 0, err
	}
	accepted := int64(0)
	for _, dev := range *devs {
		accepted += dev.Accepted
	}
	return accepted, nil
}

func (c CGMinerClient) GetTemp() (float64, error) {
	devs, err := c.Devs()
	if err != nil {
		return 0, err
	}
	// Temp is same across all boards
	return (*devs)[0].Temperature, nil
}

//func (c *BWMinerClient) Quit() error {
//	return c.client.Quit()
//}

// TODO: use the config object to decide the type of client we need
func NewClient(port ConfigPort) Client {
	c := cgminer_client.New(minerHostname, port.Get(), minerTimeout)
	return CGMinerClient{client: c}
}

var ClientModule = fx.Provide(func(port ConfigPort) Client {
	return NewClient(port)
})
