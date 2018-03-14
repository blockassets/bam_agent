package monitor_ctx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/blockassets/cgminer_client"
)

//create a mock miner to test against
type MockMiner struct {
	running bool
	stalled bool
	devs    []cgminer_client.Dev
}

func newMockMiner() *MockMiner {
	mm := &MockMiner{}
	return mm
}

func (mm *MockMiner) Start() {
	// create a new blank devs array
	mm.devs = make([]cgminer_client.Dev, 4)
	mm.running = true
}

func (mm *MockMiner) Stall() {
	mm.stalled = true
}

// funcs from service.Miner Interface
func (mm *MockMiner) Devs() (*[]cgminer_client.Dev, error) {
	if !mm.running {
		return nil, errors.New("MockMiner not running")
	}
	if !mm.stalled {
		for i, _ := range mm.devs {
			// have to index as we want to change the value
			mm.devs[i].Accepted += 1
		}
	}
	return &mm.devs, nil
}

func (mm *MockMiner) Quit() error {
	mm.running = false
	return nil
}

func TestAcceptedMonitor(t *testing.T) {

	config := &AcceptedConfig{Enabled: true, Period: time.Millisecond * 50}

	// Test for three main conditions
	// 1) Happy path: i.e. accepted shares continue to rise
	// 2) Test with stall
	// 3) Test with a miner that is not there or error
	// 4) Test for a restart of a miner between tests
	testAcceptedSharesRise(t, config)
	testAcceptedSharesStall(t, config)
	testAcceptedSharesMinerQuit(t, config)
	testAcceptedSharesMinerRestart(t, config)

}

func testAcceptedSharesRise(t *testing.T,  config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }
	// Need our own miner as the monitor tests effectively run in parallel
	mockMiner := newMockMiner()
	mockMiner.Start()


	monitors := &[]Monitor{
		NewAcceptedMonitor(config, mockMiner, onStall),
	}
	stopMonitors := StartMonitors(context.Background(), *monitors)

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period * 2)

	stopMonitors()

	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}

func testAcceptedSharesStall(t *testing.T,  config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }

	// Need our own miner as the monitor tests effectively run in parallel
	mockMiner := newMockMiner()
	mockMiner.Start()
	mockMiner.Stall()

	monitors := &[]Monitor{
		NewAcceptedMonitor(config, mockMiner, onStall),
	}
	stopMonitors := StartMonitors(context.Background(), *monitors)

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period * 2)

	stopMonitors()

	if stallCount == 0 {
		t.Errorf("Expected stallCount to be > 0")
	}
}

func testAcceptedSharesMinerQuit(t *testing.T, config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }
	// Need our own miner as the monitor tests effectively run in parallel
	mockMiner := newMockMiner()
	mockMiner.Start()
	mockMiner.Quit()

	monitors := &[]Monitor{
		NewAcceptedMonitor(config, mockMiner, onStall),
	}
	stopMonitors := StartMonitors(context.Background(), *monitors)
	// Sleep to ensure the timer is mid cycle
	time.Sleep(config.Period * 2)
	stopMonitors()

	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}

func testAcceptedSharesMinerRestart(t *testing.T, config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }
// Need our own miner as the monitor tests effectively run in parallel
	mockMiner := newMockMiner()
	mockMiner.Start()

	monitors := &[]Monitor{
		NewAcceptedMonitor(config, mockMiner, onStall),
	}
	stopMonitors := StartMonitors(context.Background(), *monitors)
	// Sleep to ensure the timer has a cycled
	time.Sleep(config.Period * 2)
	// restart the miner
	mockMiner.Quit()
	time.Sleep(config.Period * 2)
	mockMiner.Start()
	// get another cycle
	time.Sleep(config.Period * 2)

	stopMonitors()

	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}
