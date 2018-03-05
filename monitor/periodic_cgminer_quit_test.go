package monitor

import (
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestPeriodicCGMQuitMonitor_Start(t *testing.T) {
	count := 0

	config := &CGMQuitConfig{Enabled: true, Period: tool.RandomDuration{Duration: time.Duration(50) * time.Millisecond}}

	context := makeContext()
	quit := func() { count++ }

	monitor := newPeriodicCGMQuit(context, config, quit)

	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period.Duration * 2)

	// Test that stop cleans up the WaitGroup
	monitor.Stop()
	context.waitGroup.Wait()

	if count == 0 {
		t.Errorf("Expected >=1 count, got %d", count)
	}
}
