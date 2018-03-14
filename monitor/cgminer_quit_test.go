package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestPeriodicCGMQuitMonitor_Start(t *testing.T) {
	count := 0
	config := &CGMQuitConfig{Enabled: true, Period: tool.RandomDuration{Duration: time.Duration(25) * time.Millisecond}}
	quit := func() { count++ }

	monitors := &[]Monitor{
		NewPeriodicCGMQuit(config, quit),
	}
	stopMonitors := StartMonitors(context.Background(), *monitors)

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period.Duration * 2)

	// Test that stop cleans up the WaitGroup
	stopMonitors()

	if count == 0 {
		t.Errorf("Expected >=1 count, got %d", count)
	}
}
