package monitor_ctx

import "time"

// Simple PeriodicMonitor type that implements
// the Monitor interface.
//
type PeriodicMonitor struct {
	enabled     bool
	period      time.Duration
	monitorFunc MonitorFunc
}

func newPeriodicMonitor(enabled bool, period time.Duration, monitorFunc MonitorFunc) *PeriodicMonitor {
	return &PeriodicMonitor{enabled, period, monitorFunc}
}

func (pm *PeriodicMonitor) Enabled() bool {
	return pm.enabled
}

func (pm *PeriodicMonitor) Period() time.Duration {
	return pm.period
}

func (pm *PeriodicMonitor) MonitorFunc() MonitorFunc {
	return pm.monitorFunc
}
