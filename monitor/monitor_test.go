package monitor

import (
	"context"
	"testing"
	"time"
)

func NewMockMonitor(enabled bool, period time.Duration, onTick OnTick) Monitor {
	return &Data{
		Enabled: enabled,
		Period:  period * time.Millisecond,
		OnTick:  onTick,
	}
}

func TestStartMonitors(t *testing.T) {
	count1 := 0
	count2 := 0
	count3 := 0

	onTicker1 := func() TickerFunc { return func(ctx context.Context) { count1++ } }
	onTicker2 := func() TickerFunc { return func(ctx context.Context) { count2++ } }
	onTicker3 := func() TickerFunc { return func(ctx context.Context) { count3++ } }

	monitors := []Monitor{
		NewMockMonitor(true, time.Duration(10), onTicker1),
		NewMockMonitor(true, time.Duration(30), onTicker2),
		NewMockMonitor(false, time.Duration(20), onTicker3),
	}

	// Test they start and run
	stopGroup1 := StartMonitors(context.Background(), monitors)
	time.Sleep(75 * time.Millisecond)
	stopGroup1()

	if count1 < 3 {
		t.Fatalf("expected count1 to be greater than 2, got %v", count1)
	}
	if count2 < 2 {
		t.Fatalf("expected count2 to be at least 2, got %v", count2)
	}
	if count3 != 0 {
		t.Fatalf("expected count3 to be 0, got %v", count3)
	}

	if monitors[2].IsEnabled() {
		t.Fatalf("expected last monitor to not be enabled, got %v", monitors[2].IsEnabled())
	}
}

func TestStopMonitors(t *testing.T) {

	count1 := 0
	count2 := 0
	count3 := 0

	onTicker1 := func() TickerFunc { return func(ctx context.Context) { count1++ } }
	onTicker2 := func() TickerFunc { return func(ctx context.Context) { count2++ } }
	onTicker3 := func() TickerFunc { return func(ctx context.Context) { count3++ } }

	monitors := []Monitor{
		NewMockMonitor(true, time.Duration(10), onTicker1),
		NewMockMonitor(true, time.Duration(30), onTicker2),
		NewMockMonitor(false, time.Duration(20), onTicker3),
	}

	// Test they start and run
	stopGroup1 := StartMonitors(context.Background(), monitors)
	time.Sleep(18 * time.Millisecond)
	stopGroup1()
	// make sure they stop
	time.Sleep(15 * time.Millisecond)

	// Timing sensitive
	if count1 == 0 {
		t.Fatalf("expected count1 to be >=1, got %v", count1)
	}

	if count2 != 0 {
		t.Fatalf("expected count1 to be 0, got %v", count2)
	}

	if count3 != 0 {
		t.Fatalf("expected count3 to be 0, got %v", count3)
	}

	if monitors[2].IsEnabled() {
		t.Fatalf("expected last monitor to not be enabled, got %v", monitors[2].IsEnabled())
	}
}
