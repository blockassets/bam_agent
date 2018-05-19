package os

// Type insurance
var _ Ntpdate = &MockNtpdate{}

type MockNtpdate struct {
	Counter       int
	CalledNtpdate bool
}

func (ms *MockNtpdate) Ntpdate() error {
	ms.Counter++
	ms.CalledNtpdate = true
	return nil
}

func NewMockNtpdate() MockNtpdate {
	return MockNtpdate{}
}
