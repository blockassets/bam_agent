package monitor

import "time"

// Simple PeriodicMonitor type that implements
// the Monitor interface.
//
type Periodic struct {
	enabled     bool
	period      time.Duration
	monitorFunc MonitorFunc
}

func (pm *Periodic) Enabled() bool {
	return pm.enabled
}

func (pm *Periodic) Period() time.Duration {
	return pm.period
}

func (pm *Periodic) MonitorFunc() MonitorFunc {
	return pm.monitorFunc
}
