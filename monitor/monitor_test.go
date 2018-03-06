package monitor

import (
	"testing"
	"time"
)

const (
	monitorCycleDuration = time.Duration(50) * time.Millisecond
)

type TestMonitor struct {
	*Context
	CountIt func()
}

func newTestMonitor(context *Context, doItFunc func()) *TestMonitor {
	return &TestMonitor{
		Context: context,
		CountIt: doItFunc}
}

func (monitor *TestMonitor) Start() error {

	go monitor.makeTimerFunc(monitor.CountIt, monitorCycleDuration)()

	return nil
}

type testManager struct {
	Manager
}

var count int

func (mgr *testManager) StartMonitors() {
	context := makeContext()
	doIt := func() { count++ }

	mgr.Monitors = &[]Monitor{newTestMonitor(context, doIt)}

	for _, monitor := range *mgr.Monitors {
		monitor.Start()
	}
}

func TestManager_StopMonitors(t *testing.T) {
	monitorManager := &testManager{Manager{Config: nil, Client: nil}}

	monitorManager.StopMonitors()
	monitorManager.StartMonitors()
	// let enough time to go through 1 cycle
	time.Sleep(monitorCycleDuration * 2)
	monitorManager.StopMonitors()
	monitorManager.StopMonitors()
	if count != 1 {
		t.Errorf("Expected count to be 1, got %v", count)
	}
}
