package monitor

import (
	"errors"
	"log"
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
			log.Println("dev.Accepted =", mm.devs[i].Accepted)
		}
	}
	return &mm.devs, nil
}

func (mm *MockMiner) Quit() error {
	mm.running = false
	return nil
}

func TestAcceptedMonitor(t *testing.T) {
	mockMiner := newMockMiner()
	config := &AcceptedConfig{Enabled: true, Period: time.Millisecond * 50}

	// Test for three main conditions
	// 1) Happy path: i.e. accepted shares continue to rise
	// 2) Test with stall
	// 3) Test with a miner that is not there or error
	// 4) Test for a restart of a miner between tests

	mockMiner.Start()
	testAcceptedSharesRise(t, mockMiner, config)
	testAcceptedSharesStall(t, mockMiner, config)
	testAcceptedSharesMinerQuit(t, mockMiner, config)
	testAcceptedSharesMinerRestart(t, mockMiner, config)

}

func testAcceptedSharesRise(t *testing.T, mockMiner *MockMiner, config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }

	// Monitors are only meant to run once... behaviour is undefined for starting an instance twice
	context := makeContext()
	monitor := newAcceptedMonitor(context, config, mockMiner, onStall)
	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period * 2)

	monitor.Stop()
	// Make sure monitor is finished before testing results
	context.waitGroup.Wait()

	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}

func testAcceptedSharesStall(t *testing.T, mockMiner *MockMiner, config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }

	mockMiner.Stall()

	context := makeContext()
	monitor := newAcceptedMonitor(context, config, mockMiner, onStall)
	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}
	// Sleep to ensure the timer runs once
	time.Sleep(config.Period * 2)

	monitor.Stop()
	// Make sure monitor is finished before testing results
	context.waitGroup.Wait()

	if stallCount == 0 {
		t.Errorf("Expected stallCount to be > 0")
	}
}

func testAcceptedSharesMinerQuit(t *testing.T, mockMiner *MockMiner, config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }

	mockMiner.Quit()
	context := makeContext()
	monitor := newAcceptedMonitor(context, config, mockMiner, onStall)
	// Sleep to ensure the timer is mid cycle
	time.Sleep(config.Period * 2)

	monitor.Stop()
	// Make sure monitor is finished before testing results
	context.waitGroup.Wait()

	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}

func testAcceptedSharesMinerRestart(t *testing.T, mockMiner *MockMiner, config *AcceptedConfig) {
	stallCount := 0
	onStall := func() { stallCount++ }

	mockMiner.Start()
	context := makeContext()
	monitor := newAcceptedMonitor(context, config, mockMiner, onStall)
	// Sleep to ensure the timer has a cycled
	time.Sleep(config.Period * 2)
	// restart the miner
	mockMiner.Quit()
	mockMiner.Start()
	// get another cycle
	time.Sleep(config.Period * 2)

	monitor.Stop()
	// Make sure monitor is finished before testing results
	context.waitGroup.Wait()

	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}
