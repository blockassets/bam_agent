package monitor

import (
	"testing"
	"time"
)

func TestAcceptedMonitor(t *testing.T) {
	// Counters to check the closures are called by the monitor correctly
	stallCount := 0
	callCount := 0
	lastAccepted := int64(0)

	config := &AcceptedConfig{Enabled: true, Period: time.Millisecond * 50}
	context := makeContext()
	onStall := func() { stallCount++ }
	getAcceptedNoStall := func() int64 {
		callCount++
		lastAccepted += 1
		return lastAccepted
	}
	getAcceptedStall := func() int64 {
		callCount++
		return lastAccepted
	}

	// test with no stall (ie accepted shares continue to rise)
	monitor := newAcceptedMonitor(context, config, getAcceptedNoStall, onStall)
	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period * 2)

	monitor.Stop()

	if callCount == 0 {
		t.Errorf("Expected >=1 callCount, got %d", callCount)
	}
	if stallCount != 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}

	// Setup and Test with stall
	stallCount = 0
	context = makeContext()
	monitor = newAcceptedMonitor(context, config, getAcceptedStall, onStall)
	err = monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(config.Period * 2)

	monitor.Stop()

	if stallCount == 0 {
		t.Errorf("Expected stallCount to be 0, got %d", stallCount)
	}
}
