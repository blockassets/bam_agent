package monitor_ctx

import (
	"context"
	"testing"
	"time"
)

func TestStartMonitors(t *testing.T) {
	count1 := 0
	count2 := 0
	count3 := 0
	onTicker1 := func(ctx context.Context) { count1++ }
	onTicker2 := func(ctx context.Context) { count2++ }
	onTicker3 := func(ctx context.Context) { count3++ }
	monitors := &[]Monitor{
		newPeriodic(true, 50*time.Millisecond, onTicker1),
		newPeriodic(true, 75*time.Millisecond, onTicker2),
		newPeriodic(false, 100*time.Millisecond, onTicker3),
	}

	// Test they satrt and run
	stopGroup1 := StartMonitors(context.Background(), *monitors)
	time.Sleep(202 * time.Millisecond)
	stopGroup1()
	mark1 := count1
	mark2 := count2
	if mark1 < 3 {
		t.Errorf("Expected count1 to be greater than 2, got %v", mark1)
	}
	if mark2 < 2 {
		t.Errorf("Expected count2 to be at least 2, got %v", mark2)
	}
	if count3 != 0 {
		t.Errorf("Expected count3 to be 0, got %v", count3)
	}

	// make sure they stop
	time.Sleep(202 * time.Millisecond)
	if mark1 != count1 {
		t.Errorf("Expected count1 (%v) to be same as mark1(%v)", count1, mark1)
	}
	if mark2 != count2 {
		t.Errorf("Expected count2 (%v) to be same as mark2(%v)", count2, mark2)
	}
}
