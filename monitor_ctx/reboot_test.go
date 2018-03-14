package monitor_ctx

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestPeriodicRebootMonitor_Start(t *testing.T) {
	count := 0

	config := &RebootConfig{Enabled: true, Period: tool.RandomDuration{Duration: time.Duration(25) * time.Millisecond}}
	reboot := func() { count++ }

	monitors := &[]Monitor{
		NewPeriodicReboot(config, reboot),
	}
	stopMonitors := StartMonitors(context.Background(), *monitors)

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period.Duration * 2)

	stopMonitors()
	if count == 0 {
		t.Errorf("Expected >=1 count, got %d", count)
	}
}
