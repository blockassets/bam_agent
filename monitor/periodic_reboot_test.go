package monitor

import (
	"fmt"
	"testing"
	"time"
)

func TestPeriodicReboot(t *testing.T) {
	count := 0

	fmt.Printf("Starting Periodic Reboot logic test\n")
	cfg := Config{}
	cfg.Reboot = RebootConfig{Enabled: true, PeriodInSeconds: 1, InitialPeriodRangeInSeconds: 1}

	pr := newPeriodicReboot(func() { count++ })

	fmt.Printf("Starting Monitor\n")

	err := pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.1 Expected start to suceed. Returned %+v", err)
	}
	if !pr.IsRunning() {
		t.Errorf("t2.2 Expected pr.isRunning to be true")
	}
	// give it time for one call max time is 2 seconds
	time.Sleep(time.Duration(2200) * time.Millisecond)
	fmt.Printf("Stopping Monitor\n")
	pr.Stop()
	if count != 1 {
		t.Errorf("t2.3 Expected 1 on count, got %d", count)
	}
	fmt.Printf("Starting Monitor\n")

	err = pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.5 Expected 2nd start to suceced. Returned %+v", err)
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
	cfg.Reboot.Enabled = false
	fmt.Printf("Starting Monitor\n")

	err = pr.Start(&cfg)
	fmt.Printf("Started Monitor\n")

	if err != nil {
		t.Errorf("t2.8 Expected 4th start to succeed")
	}
	if !pr.IsRunning() {
		t.Errorf("t2.9 Expected to be running")
	}
	fmt.Printf("Stopping Monitor\n")
	pr.Stop()
	fmt.Printf("test getRunning\n")

	if pr.IsRunning() {
		t.Errorf("t2.10 Expected to be not running")
	}

}
