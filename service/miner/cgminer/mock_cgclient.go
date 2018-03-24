package cgminer

import "github.com/blockassets/cgminer_client"

const (
	testTemp     = 100.00
	testAccepted = 5
)

// Type insurance
var _ CgClient = &MockCgClient{}

type MockCgClient struct {
	quitCalled bool
	devs       []cgminer_client.Dev
}

func (cgClient *MockCgClient) Quit() error {
	cgClient.quitCalled = true
	return nil
}

func (cgClient *MockCgClient) Restart() error {
	return nil
}

func (cgClient *MockCgClient) Summary() (*cgminer_client.Summary, error) {
	return nil, nil
}

func (cgClient *MockCgClient) ChipStat() (*[]cgminer_client.ChipStat, error) {
	return nil, nil
}

func (cgClient *MockCgClient) Devs() (*[]cgminer_client.Dev, error) {
	return &cgClient.devs, nil
}

func newMockCgClient() *MockCgClient {
	devs := make([]cgminer_client.Dev, 1)
	devs[0].Accepted = testAccepted
	devs[0].Temperature = testTemp
	return &MockCgClient{devs: devs}
}
