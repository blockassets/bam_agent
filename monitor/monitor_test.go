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
		t.Fatalf("Expected count1 to be greater than 2, got %v", count1)
	}
	if count2 < 2 {
		t.Fatalf("Expected count2 to be at least 2, got %v", count2)
	}
	if count3 != 0 {
		t.Fatalf("Expected count3 to be 0, got %v", count3)
	}
}

func TestStopMonitors(t *testing.T) {

	//// make sure they stop
	//time.Sleep(75 * time.Millisecond)
	//if mark1 != count1 {
	//	t.Fatalf("Expected count1 (%v) to be same as mark1(%v)", count1, mark1)
	//}
	//if mark2 != count2 {
	//	t.Fatalf("Expected count2 (%v) to be same as mark2(%v)", count2, mark2)
	//}
	//
	//// they may be all disabled.
	//// Make sure that the subsystem always handles gracefully
	//monitors = []*Data{
	//	{Enabled: false, Period: time.Duration(10) * time.Millisecond, OnTick: onTicker1},
	//	{Enabled: false, Period: time.Duration(10) * time.Millisecond, OnTick: onTicker2},
	//	{Enabled: false, Period: time.Duration(10) * time.Millisecond, OnTick: onTicker3},
	//}
	//stopGroup2 := StartMonitors(context.Background(), monitors)
	//stopGroup2()
	//// And no panic on an empty array...
	//monitors = []*Monitor{}
	//stopGroup3 := StartMonitors(context.Background(), monitors)
	//stopGroup3()

}
