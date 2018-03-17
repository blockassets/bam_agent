package os

// Type insurance
var _ Miner = &MockMiner{}

type MockMiner struct {
	CalledStartMiner bool
	CalledStopMiner  bool
}

func (ms *MockMiner) Start() error {
	ms.CalledStartMiner = true
	return nil
}

func (ms *MockMiner) Stop() error {
	ms.CalledStopMiner = true
	return nil
}

func NewMockMiner() MockMiner {
	return MockMiner{}
}
