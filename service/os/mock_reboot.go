package os

// Type insurance
var _ Reboot = &MockReboot{}

type MockReboot struct {
	Counter      int
	CalledReboot bool
}

func (ms *MockReboot) Reboot() error {
	ms.Counter++
	ms.CalledReboot = true
	return nil
}

func NewMockReboot() MockReboot {
	return MockReboot{}
}
