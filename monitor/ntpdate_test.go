package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

func TestNewNtpdateMonitor(t *testing.T) {
	config := NtpdateConfig{Enabled: true, Period: time.Duration(5)}

	ntpdate := os.NewMockNtpdate()
	result := NewNtpdateMonitor(config, &ntpdate)
	monitor := result.Monitor

	if monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected period 5, got %s", monitor.GetPeriod())
	}

	if !monitor.IsEnabled() {
		t.Fatalf("expected enabled, got %v", monitor.IsEnabled())
	}

	// Run
	result.Monitor.NewTickerFunc()(context.TODO())

	if ntpdate.Counter != 1 {
		t.Fatalf("expected counter 1, got %d", ntpdate.Counter)
	}
}
