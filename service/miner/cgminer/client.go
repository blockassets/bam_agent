package cgminer

import (
	"github.com/blockassets/cgminer_client"
)

// create an interface to the actual miner client, so we can mock it out in testing
type Client interface {
	Devs() (*[]cgminer_client.Dev, error)
	Quit() error
	Restart() error
	Summary() (*cgminer_client.Summary, error)
	ChipStat() (*[]cgminer_client.ChipStat, error)
}

type ClientData struct {
	Client Client
}

func (c ClientData) Quit() error {
	return c.Client.Quit()
}

func (c *ClientData) GetAccepted() (int64, error) {
	devs, err := c.Client.Devs()
	if err != nil {
		return 0, err
	}
	accepted := int64(0)
	for _, dev := range *devs {
		accepted += dev.Accepted
	}
	return accepted, nil
}

func (c *ClientData) GetTemp() (float64, error) {
	devs, err := c.Client.Devs()
	if err != nil {
		return 0, err
	}
	// Temp is same across all boards
	return (*devs)[0].Temperature, nil
}

