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
	CalledRestart  bool
	CalledDevs     bool
	CalledAccepted bool

	test     int
	accepted int64
}

func (c *MockMinerClient) Quit() error {
	c.CalledQuit = true
	return nil
}

func (c *MockMinerClient) Restart() error {
	c.CalledRestart = true
	return nil
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
		return 0, errors.New("some error")
	}

	return c.accepted, nil
}

func (c *MockMinerClient) GetTemp() (float64, error) {
	switch c.test {
	case Under100Temp:
		return 90.0, nil
	case Exactly100Temp:
		return 100.0, nil
	case Over100Temp:
		return 101.0, nil
	}
	return -1, nil
}

func NewMockMinerClient(test int) MockMinerClient {
	return MockMinerClient{
		test:     test,
		accepted: 0,
	}
}
