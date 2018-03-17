package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

func TestNewLoadMonitor(t *testing.T) {
	config := HighLoadConfig{Enabled: true, Period: time.Duration(5), HighLoadMark: 5.0}

	reboot := os.NewMockReboot()
	sr := os.NewMockStatRetriever(os.LevelNotEnough)
	result := NewLoadMonitor(config, &sr, &reboot)
	monitor := result.Monitor

	if monitor.GetPeriod() != time.Duration(5) {
		t.Fatalf("expected period 5, got %s", monitor.GetPeriod())
	}

	if !monitor.IsEnabled() {
		t.Fatalf("expected enabled, got %v", monitor.IsEnabled())
	}

	// Run
	monitor.NewTickerFunc()(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	sr = os.NewMockStatRetriever(os.LevelBelowFive)
	result = NewLoadMonitor(config, &sr, &reboot)
	monitor = result.Monitor
	monitor.NewTickerFunc()(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	sr = os.NewMockStatRetriever(os.LevelExactlyFive)
	result = NewLoadMonitor(config, &sr, &reboot)
	monitor = result.Monitor
	monitor.NewTickerFunc()(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	sr = os.NewMockStatRetriever(os.LevelAboveFive)
	result = NewLoadMonitor(config, &sr, &reboot)
	monitor = result.Monitor
	monitor.NewTickerFunc()(context.TODO())
	if !reboot.CalledReboot {
		t.Fatalf("expected reboot, got %v", reboot.CalledReboot)
	}

	reboot = os.NewMockReboot()
	sr = os.NewMockStatRetriever(os.LevelMalformed)
	result = NewLoadMonitor(config, &sr, &reboot)
	monitor = result.Monitor
	monitor.NewTickerFunc()(context.TODO())
	if reboot.CalledReboot {
		t.Fatalf("expected no reboot, got %v", reboot.CalledReboot)
	}
}
