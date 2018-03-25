package cgminer

import "github.com/blockassets/cgminer_client"

const (
	testTemp     = 100.00
	testAccepted = 5
)

// Type insurance
var _ cgminer_client.Client = &MockCgClient{}

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
	return &cgminer_client.Summary{}, nil
}

func (cgClient *MockCgClient) ChipStat() (*[]cgminer_client.ChipStat, error) {
	cs1 := cgminer_client.ChipStat{}
	chipStats := &[]cgminer_client.ChipStat{cs1}
	return chipStats, nil
}

func (cgClient *MockCgClient) Devs() (*[]cgminer_client.Dev, error) {
	return &cgClient.devs, nil
}

func NewMockCgClient() *MockCgClient {
	devs := make([]cgminer_client.Dev, 1)
	devs[0].Accepted = testAccepted
	devs[0].Temperature = testTemp
	devs[0].Status = "Alive"
	devs[0].Enabled = "Y"

	return &MockCgClient{devs: devs}
}
