package monitor

import (
	"context"
	"sync"

	"go.uber.org/fx"
)

type Result struct {
	fx.Out
	Monitor Monitor `group:"monitor"`
}

type Group struct {
	fx.In
	Monitors []Monitor `group:"monitor"`
}

type Manager interface {
	Start()
	Stop()
}

type ManagerData struct {
	monitors      []Monitor
	startCount    int
	stopFunc      func()
	startMonitors func(context.Context, []Monitor) Stop
	sync.Mutex
}

func (mm *ManagerData) Start() {
	mm.Lock()
	defer mm.Unlock()
	if mm.startCount == 0 {
		mm.stopFunc = mm.startMonitors(context.Background(), mm.monitors)
	}
	mm.startCount++
}

func (mm *ManagerData) Stop() {
	mm.Lock()
	defer mm.Unlock()
	if mm.startCount == 1 {
		mm.stopFunc()
		mm.stopFunc = nil
	}
	mm.startCount--
}

func StartMonitors(parent context.Context, monitors []Monitor) Stop {
	ctx, cancelFunc := context.WithCancel(parent)
	wg := sync.WaitGroup{}
	for _, mon := range monitors {
		m := mon
		if m.IsEnabled() {
			m.Start(ctx, &wg, m.NewTickerFunc())
		}
	}
	return func() {
		cancelFunc()
		wg.Wait()
	}
}

/*
	There is a little fx magic here. Group gets magically populated with a list of
	monitors because those are provided in the module declaration below and we
	use the fx 'group' functionality to make that happen.
*/
func NewManager(g Group) Manager {
	return &ManagerData{
		monitors:      g.Monitors,
		startMonitors: StartMonitors,
	}
}

var Module = fx.Options(
	ConfigProviders,

	fx.Provide(
		NewManager,

		NewAcceptedMonitor,
		NewCGMQuitMonitor,
		NewHighTempMonitor,
		NewLowMemoryMonitor,
		NewLoadMonitor,
		NewNtpdateMonitor,
		NewRebootMonitor,
	),
)
