package monitor

import (
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestPeriodicRebootMonitor_Start(t *testing.T) {
	count := 0

	context := makeContext()
	config := &RebootConfig{Enabled: true, Period: tool.RandomDuration{Duration: time.Duration(50) * time.Millisecond}}
	reboot := func() { count++ }

	monitor := newPeriodicReboot(context, config, reboot)

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
