package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/miner"
)

func TestNewCGMQuitMonitor(t *testing.T) {
	config := CGMQuitConfig{Enabled: true, Period: time.Duration(5)}

	client := miner.NewMockMinerClient(-1)
	result := NewCGMQuitMonitor(config, &client)
	monitor := result.Monitor

	if monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected period 5, got %s", monitor.GetPeriod())
	}

	if !monitor.IsEnabled() {
		t.Fatalf("expected enabled, got %v", monitor.IsEnabled())
	}

	// Run
	result.Monitor.NewTickerFunc()(context.TODO())

	if !client.CalledQuit {
		t.Fatalf("expected client.calledQuit true, got %v", client.CalledQuit)
	}
}
