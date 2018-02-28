package monitor

import (
	"fmt"
	"testing"
	"time"
)

func TestPeriodicReboot(t *testing.T) {
	count := 0

	fmt.Printf("Starting Periodic Reboot logic test")
	cfg := MonitorConfig{}
	cfg.Reboot = RebootConfig{Enabled: true, PeriodSecs: 1, InitialPeriodRange: 1}

	pr := newPeriodicReboot(func() { count++ })

	fmt.Printf("Starting Monitor")

	err := pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.1 Expected start to suceed. Returned %+v", err)
	}
	if pr.getRunning() != true {
		t.Errorf("t2.2 Expected pr.isRunning to be true")
	}
	// give it time for one call max tim eis 2 seconds
	time.Sleep(time.Duration(1500) * time.Millisecond)
	fmt.Printf("Stopping Monitor")
	pr.Stop()
	if count != 1 {
		t.Errorf("t2.3 Expected 1 on count, got %d", count)
	}
	fmt.Printf("Starting Monitor")

	err = pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.5 Expected 2nd start to suceed. Returned %+v", err)
	}
	fmt.Printf("Starting Monitor")

	err = pr.Start(&cfg)
	if err == nil {
		t.Errorf("t2.6 Expected 3rd start to fail")
	}
	fmt.Printf("Stopping Monitor")

	pr.Stop()
	if pr.getRunning() {
		t.Errorf("t2.7 Expected to be not running")
	}
	cfg.Load.Enabled = false
	fmt.Printf("Starting Monitor")

	err = pr.Start(&cfg)
	if err != nil {
		t.Errorf("t2.8 Expected 4th start to succeed")
	}
	if !pr.getRunning() {
		t.Errorf("t2.9 Expected to be running")
	}
	fmt.Printf("Stopping Monitor")

	pr.Stop()
	if pr.getRunning() {
		t.Errorf("t2.10 Expected to be not running")
	}

}
