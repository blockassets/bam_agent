package monitor

import (
	"testing"
	"time"
)

func TestPeriodicRebootMonitor_Start(t *testing.T) {
	count := 0

	context := makeContext()
	config := &RebootConfig{Enabled: true, Period: 1}
	initialPeriod := time.Duration(50) * time.Millisecond
	reboot := func() { count++ }

	monitor := newPeriodicReboot(context, config, initialPeriod, reboot)

	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(initialPeriod * 2)

	// Test that stop cleans up the WaitGroup
	monitor.Stop()
	context.waitGroup.Wait()

	if count == 0 {
		t.Errorf("Expected >=1 count, got %d", count)
	}
}
