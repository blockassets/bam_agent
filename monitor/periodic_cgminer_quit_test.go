package monitor

import (
	"fmt"
	"testing"
	"time"
)

func TestPeriodicCGMQuit(t *testing.T) {
	count := 0

	fmt.Printf("Starting Periodic CGMQuit logic test\n")
	cfg := Config{}
	cfg.CGMQuit = CGMQuitConfig{Enabled: true, PeriodInSeconds: 1, InitialPeriodRangeInSeconds: 1}

	pr := newPeriodicCGMQuit(func() { count++ })

	fmt.Printf("Starting Monitor\n")

	err := pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.1 Expected start to succeed. Returned %+v", err)
	}
	if !pr.IsRunning() {
		t.Errorf("t2.2 Expected pr.isRunning to be true")
	}

	// give it time for an inital call(between 1 and 2 seconds) and then one more..
	time.Sleep(time.Duration(3200) * time.Millisecond)

	fmt.Printf("Stopping Monitor")
	pr.Stop()
	mark := count
	if (count < 2) || (count > 3) {
		t.Errorf("t2.3 Expected 2 or 3 on count, got %d", count)
	}
	time.Sleep(time.Duration(3) * time.Second)
	if count != mark {
		t.Errorf("t2.4 Expected count to be the same as mark, %d != %d", count, mark)
	}
	fmt.Printf("Starting Monitor\n")
	err = pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.5 Expected 2nd start to succeed. Returned %+v", err)
	}

	fmt.Printf("Starting Monitor\n")
	err = pr.Start(&cfg)
	if err == nil {
		t.Errorf("t2.6 Expected 3rd start to fail")
	}

	fmt.Printf("Stopping Monitor\n")
	pr.Stop()
	if pr.IsRunning() {
		t.Errorf("t2.7 Expected to be not running")
	}

	fmt.Printf("Starting Monitor\n")
	cfg.CGMQuit.Enabled = false
	err = pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.8 Expected 4th start to succeed")
	}
	if !pr.IsRunning() {
		t.Errorf("t2.9 Expected to be running")
	}

	fmt.Printf("Stopping Monitor\n")
	pr.Stop()
	if pr.IsRunning() {
		t.Errorf("t2.10 Expected to be not running")
	}

	count = 0
	cfg.CGMQuit.Enabled = true
	err = pr.Start(&cfg)
	time.Sleep(time.Duration(2500) * time.Millisecond)
	if count < 1 {
		t.Errorf("t2.10 Expected count to be non-zero, got %d", count)
	}

}
