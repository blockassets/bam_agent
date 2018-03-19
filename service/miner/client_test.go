package miner

import (
	"testing"

	"github.com/blockassets/bam_agent/tool"
	"github.com/blockassets/cgminer_client"
)

func TestNewClient(t *testing.T) {
	config := NewConfig(tool.NewCmdLine())
	client := NewClient(NewConfigPort(config))
	err := client.Quit()
	if err == nil {
		t.Fatal("expected an error when calling quit!")
	}
}

// Type insurance
var _ CgminerClientInterface = &cgminer_client.Client{}

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

const (
	testTemp     = 100.00
	testAccepted = 5
)

func makeMockCgClient() *MockCgClient {
	devs := make([]cgminer_client.Dev, 1)
	devs[0].Accepted = testAccepted
	devs[0].Temperature = testTemp
	return &MockCgClient{devs: devs}
}

type MockCgClient struct {
	quitCalled bool
	devs       []cgminer_client.Dev
}

func TestCGMinerClient_Quit(t *testing.T) {
	cgclient := makeMockCgClient()
	client := &CGMinerClient{client: cgclient}

	err := client.Quit()
	if err != nil {
		t.Fatal("Did not expected an error when calling quit!")
	}
	if !cgclient.quitCalled {
		t.Fatal("Expected quitCalled to be true!")
	}
}

func TestCGMinerClient_GetAccepted(t *testing.T) {
	cgclient := makeMockCgClient()
	client := &CGMinerClient{client: cgclient}

	accepted, err := client.GetAccepted()
	if err != nil {
		t.Fatal("Did not expected an error when calling GetAccepted!")
	}
	if accepted != testAccepted {
		t.Fatalf("Expected accepted to be equal to testAccepted, got %v", accepted)
	}
}

func TestCGMinerClient_GetTemp(t *testing.T) {
	cgclient := makeMockCgClient()
	client := &CGMinerClient{client: cgclient}

	temp, err := client.GetTemp()
	if err != nil {
		t.Fatal("Did not expected an error when calling GetTemp!")
	}
	if temp != testTemp {
		t.Fatalf("Expected accepted to be equal to testTemp, got %v", temp)
	}
}
