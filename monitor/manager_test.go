package monitor

import (
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	count := 0
	onLoadHigh := func() { count++ }
	config := &HighLoadConfig{Enabled: true, Period: 25 * time.Millisecond, HighLoadMark: 5.0}
	sr := NewTestStatRetriever(LevelAboveFive)

	mm := &Manager{}
	mm.monitors = &[]Monitor{
		NewLoadMonitor(config, sr, onLoadHigh),
	}
	mm.Start()
	time.Sleep(config.Period * 2)
	mm.Stop()
	mark := count
	if mark == 0 {
		t.Errorf("Expected >=1 count, got %d", mark)
	}
	time.Sleep(config.Period * 2)
	if mark != count {
		t.Errorf("Expected no change in count(mark = %v), got %v", mark, count)
	}

	testNestedStopStarts(t, mm, &count, config.Period)
}

func testNestedStopStarts(t *testing.T, mm *Manager, count *int, period time.Duration) {
	// Simulate Requeat A stoppingMonitors... Request B stopping... Request A starting... Request B starting
	// Expected results are that monitors dont run from Request A until Request B starts
	*count = 0
	mark := 0
	mm.Start() // simulates init
	time.Sleep(period * 2)
	if *count == 0 {
		t.Errorf("Expected >0 *count ")
	}
	mm.Stop() // Request A
	mark = *count
	time.Sleep(period * 2)
	// test we were stopped
	if mark != *count {
		t.Errorf("Expected no change in *count(mark = %v), got %v", mark, *count)
	}
	mm.Stop() // Request B
	mark = *count
	time.Sleep(period * 2)
	// test we were stopped
	if mark != *count {
		t.Errorf("Expected no change in *count(mark = %v), got %v", mark, *count)
	}
	mm.Start() // Request A or B starts
	mark = *count
	time.Sleep(period * 2)
	// We should STILL be stopped as we are nesting...
	if mark != *count {
		t.Errorf("Expected no change in *count(mark = %v), got %v", mark, *count)
	}
	*count = 0
	mm.Start() // Request B  or A starts
	time.Sleep(period * 2)
	// Expect to start for realz...
	if *count == 0 {
		t.Errorf("Expected increase in *count ")
	}
}
