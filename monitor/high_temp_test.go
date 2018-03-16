package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
)

func TestHighTempMonitor(t *testing.T) {
	config := HighTempConfig{Enabled: true, Period: time.Duration(5), HighTemp: 100}

	mockMinerClient := miner.NewMockMinerClient(miner.Under100Temp)
	mockMinerOs := os.NewMockMiner()
	result := NewHighTempMonitor(config, &mockMinerClient, &mockMinerOs)

	if !result.Monitor.IsEnabled() {
		t.Fatalf("expected monitor enabled, got %v", result.Monitor.IsEnabled())
	}

	if result.Monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected monitor period 5, got %s", result.Monitor.GetPeriod())
	}

	// Run
	result.Monitor.NewTickerFunc()(context.TODO())
	if mockMinerOs.CalledStartMiner {
		t.Fatalf("expected calledStopMiner to be false, got %v", mockMinerOs.CalledStopMiner)
	}

	mockMinerClient = miner.NewMockMinerClient(miner.Exactly100Temp)
	mockMinerOs = os.NewMockMiner()
	result = NewHighTempMonitor(config, &mockMinerClient, &mockMinerOs)
	result.Monitor.NewTickerFunc()(context.TODO())

	if !mockMinerOs.CalledStopMiner {
		t.Fatalf("expected calledStopMiner to be true, got %v", mockMinerOs.CalledStopMiner)
	}

	mockMinerClient = miner.NewMockMinerClient(miner.Over100Temp)
	mockMinerOs = os.NewMockMiner()
	result = NewHighTempMonitor(config, &mockMinerClient, &mockMinerOs)
	result.Monitor.NewTickerFunc()(context.TODO())

	if !mockMinerOs.CalledStopMiner {
		t.Fatalf("expected calledStopMiner to be true, got %v", mockMinerOs.CalledStopMiner)
	}
}
