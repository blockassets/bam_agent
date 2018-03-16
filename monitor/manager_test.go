package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/os"
)

func TestManager(t *testing.T) {
	config := HighLoadConfig{Enabled: true, Period: 25 * time.Millisecond, HighLoadMark: 5.0}
	sr := os.NewMockStatRetriever(os.LevelAboveFive)
	reboot := os.NewMockReboot()
	result := NewLoadMonitor(config, &sr, &reboot)

	mm := NewManager(Group{
		Monitors: []Monitor{
			result.Monitor,
		},
	})

	mm.Start()
	time.Sleep(config.Period * 2)
	mm.Stop()

	// By the time we stop, we have run the reboot at least once
	if reboot.Counter == 0 {
		t.Fatalf("Expected == 0 count, got %d", reboot.Counter)
	}
}

func TestMultipleStopStarts(t *testing.T) {
	stopCount := 0
	mm := &ManagerData{
		startMonitors: func(context context.Context, monitors []Monitor) Stop {
			stopCount = 0
			return func() {
				stopCount++
			}
		},
	}

	// Just a normal start/stop
	mm.Start()
	if mm.startCount != 1 {
		t.Fatalf("expected startCount 1, got %v", mm.startCount)
	}
	mm.Stop()
	if mm.startCount != 0 {
		t.Fatalf("expected startCount 0, got %v", mm.startCount)
	}
	if stopCount != 1 {
		t.Fatalf("expected stopCount 1, got %v", stopCount)
	}

	// Does not actually stop the monitors a second time
	mm.Stop()
	if stopCount != 1 {
		t.Fatalf("expected stopCount 1, got %v", stopCount)
	}
	if mm.startCount != -1 {
		t.Fatalf("expected startCount -1, got %v", mm.startCount)
	}

	// This should not do anything except get us back to startCount 0
	mm.Start()
	if stopCount != 1 {
		t.Fatalf("expected stopCount 1, got %v", stopCount)
	}
	if mm.startCount != 0 {
		t.Fatalf("expected startCount 0, got %v", mm.startCount)
	}

	// This should reset us back to 0 on stopCount since we are starting fresh again
	mm.Start()
	if stopCount != 0 {
		t.Fatalf("expected stopCount 0, got %v", stopCount)
	}
	if mm.startCount != 1 {
		t.Fatalf("expected startCount 1, got %v", mm.startCount)
	}
}
