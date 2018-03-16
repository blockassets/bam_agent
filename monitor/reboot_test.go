package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

func TestNewRebootMonitor(t *testing.T) {
	config := RebootConfig{Enabled: true, Period: time.Duration(5)}

	reboot := os.NewMockReboot()
	result := NewRebootMonitor(config, &reboot)
	monitor := result.Monitor

	if monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected period 5, got %s", monitor.GetPeriod())
	}

	if !monitor.IsEnabled() {
		t.Fatalf("expected enabled, got %v", monitor.IsEnabled())
	}

	// Run
	result.Monitor.NewTickerFunc()(context.TODO())

	if reboot.Counter != 1 {
		t.Fatalf("expected counter 1, got %d", reboot.Counter)
	}
}
