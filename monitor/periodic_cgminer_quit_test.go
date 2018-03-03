package monitor

import (
	"sync"
	"testing"
	"time"
)

func TestPeriodicCGMQuitMonitor_Start(t *testing.T) {
	count := 0

	config := &CGMQuitConfig{Enabled: true, Period: 1}
	context := &Context{quit: make(chan bool), waitGroup: &sync.WaitGroup{}}
	initialPeriod := time.Duration(50) * time.Millisecond
	quit := func() { count++ }

	monitor := newPeriodicCGMQuit(context, config, &initialPeriod, quit)

	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(initialPeriod * 2)

	// Test that stop cleans up the WaitGroup
	monitor.Stop()
	context.waitGroup.Wait()

	if count == 0 {
		t.Errorf("Expected >=1 count, got %d", count)
	}
}
