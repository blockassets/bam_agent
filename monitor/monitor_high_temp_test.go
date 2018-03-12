package monitor

import (
	"errors"
	"testing"
	"time"

	"github.com/blockassets/cgminer_client"
)

const (
	Under100 = iota
	Exactly100
	Over100
	ReturnError
)

//create a mock miner to test against
type TempMockMiner struct {
	test int
	devs []cgminer_client.Dev
}

func newTempMockMiner() *TempMockMiner {
	mm := &TempMockMiner{}
	return mm
}

func (mm *TempMockMiner) setTest(test int) {
	mm.test = test
	return
}

func (mm *TempMockMiner) Start() {
	// create a new blank devs array
	mm.devs = make([]cgminer_client.Dev, 4)
}

func (mm *TempMockMiner) Quit() error {
	return nil
}

// funcs from service.Miner Interface
func (mm *TempMockMiner) Devs() (*[]cgminer_client.Dev, error) {
	for i, _ := range mm.devs {
		// have to index as we want to change the value
		mm.devs[i].Temperature = 99
	}
	// vary temp of one of the devices depending on test
	switch mm.test {
	case Under100:
		mm.devs[0].Temperature = 90.0
	case Exactly100:
		mm.devs[0].Temperature = 100.0
	case Over100:
		mm.devs[0].Temperature = 101.0
	}
	if mm.test == ReturnError {
		return nil, errors.New("highTempTest")
	}

	return &mm.devs, nil
}

func TestHighTempMonitor(t *testing.T) {
	tempMockMiner := newTempMockMiner()
	config := &HighTempConfig{Enabled: true, Period: time.Millisecond * 50, HighTemp: 100}
	context := makeContext()
	highTempCount := 0
	onHighTemp := func() { highTempCount++ }
	monitor := newHighTempMonitor(context, config, tempMockMiner, onHighTemp)

	tempMockMiner.Start()
	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	tempMockMiner.setTest(Under100)
	// Sleep to ensure the timer runs once
	highTempCount = 0
	time.Sleep(config.Period * 2)
	if highTempCount != 0 {
		t.Errorf("Expected highTempCount to be 0, got %d", highTempCount)
	}

	tempMockMiner.setTest(Exactly100)
	// Sleep to ensure the timer runs once
	highTempCount = 0
	time.Sleep(config.Period * 2)
	if highTempCount == 0 {
		t.Errorf("Expected highTempCount to be greater than 0")
	}

	tempMockMiner.setTest(Over100)
	// Sleep to ensure the timer runs once
	highTempCount = 0
	time.Sleep(config.Period * 2)
	if highTempCount == 0 {
		t.Errorf("Expected highTempCount to be greater than 0")
	}
	tempMockMiner.setTest(ReturnError)
	// Sleep to ensure the timer runs once
	highTempCount = 0
	time.Sleep(config.Period * 2)
	if highTempCount != 0 {
		t.Errorf("Expected highTempCount to be 0, got %d", highTempCount)
	}

	// Test that stop cleans up the WaitGroup
	monitor.Stop()
	context.waitGroup.Wait()
}
