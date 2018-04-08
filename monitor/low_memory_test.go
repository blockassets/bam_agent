package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

func TestNewLowMemoryMonitor(t *testing.T) {
	config := LowMemoryConfig{Enabled: true, Period: time.Duration(5), LowMemory: 125 * 1000 * 1000}

	reboot := os.NewMockReboot()
	sr := os.NewMockMemInfo(os.MemBelow)
	result := NewLowMemoryMonitor(config, &sr, &reboot)
	monitor := result.Monitor

	if monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected period 5, got %s", monitor.GetPeriod())
	}

	if !monitor.IsEnabled() {
		t.Fatalf("expected enabled, got %v", monitor.IsEnabled())
	}

	// Run
	monitor.NewTickerFunc()(context.TODO())
	if !reboot.CalledReboot {
		t.Fatalf("expected reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	sr = os.NewMockMemInfo(os.MemExactly)
	result = NewLowMemoryMonitor(config, &sr, &reboot)
	monitor = result.Monitor
	monitor.NewTickerFunc()(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	sr = os.NewMockMemInfo(os.MemAbove)
	result = NewLowMemoryMonitor(config, &sr, &reboot)
	monitor = result.Monitor
	monitor.NewTickerFunc()(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}
}
