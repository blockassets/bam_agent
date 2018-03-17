package miner

import (
	"errors"
)

const (
	Under100Temp = iota
	Exactly100Temp
	Over100Temp
	AcceptedIncrement
	AcceptedSame
	AcceptedZero
	AcceptedError
)

// Type insurance
var _ Client = &MockMinerClient{}

type MockMinerClient struct {
	CalledQuit     bool
	CalledDevs     bool
	CalledAccepted bool

	test     int
	accepted int64
	devs     []Dev
}

func (c *MockMinerClient) Quit() error {
	c.CalledQuit = true
	return nil
}

func (c *MockMinerClient) Devs() (*[]Dev, error) {
	c.CalledDevs = true
	switch c.test {
	case Under100Temp:
		c.devs[0].Temperature = 90.0
	case Exactly100Temp:
		c.devs[0].Temperature = 100.0
	case Over100Temp:
		c.devs[0].Temperature = 101.0
	}
	return &c.devs, nil
}

func (c *MockMinerClient) GetAccepted() (int64, error) {
	c.CalledAccepted = true
	switch c.test {
	case AcceptedIncrement:
		c.accepted++
	case AcceptedSame:
		c.accepted = 1
	case AcceptedZero:
		c.accepted = 0
	case AcceptedError:
		c.accepted = -1
	}

	if c.accepted == -1 {
		return -1, errors.New("some error")
	}

	return c.accepted, nil
}

func (c *MockMinerClient) GetTemp() (float64, error) {
	devs, err := c.Devs()
	if err != nil {
		return -1, err
	}
	dev := (*devs)[0]
	return dev.Temperature, nil
}

func NewMockMinerClient(test int) MockMinerClient {
	return MockMinerClient{
		test:     test,
		accepted: 0,
		devs:     make([]Dev, 1),
	}
}
