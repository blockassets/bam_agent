package monitor_ctx

import (
	"context"
	"sync"
	"time"
)

// Monitors need to be able to be safely cancelled as a group and they should each be able to
// be able clean themselves up on a cancel.
//
// This intends to create immutable monitors. Once they are started only the returned stop function can stop them.
//
// Monitors are started in groups and the StartMonitors function returns a function to stop them all
// This function will block until all the monitors have finished cleaning up.
//
// If this function is not called, the monitors associated with it will continue on until the
// application terminates. Uses the context package to implement
//
// see https://golang.org/pkg/context/#CancelFunc and
//     https://blog.golang.org/context

type MonitorFunc func(ctx context.Context)

type Monitor interface {
	Enabled() bool
	Period() time.Duration
	MonitorFunc() MonitorFunc
}

func StartMonitors(parent context.Context, monitors []Monitor) (stopFunc func()) {
	ctx, cancelFunc := context.WithCancel(parent)
	wg := sync.WaitGroup{}
	for _, mon := range monitors {
		if mon.Enabled() {
			runMonitor(ctx, &wg, mon.Period(), mon.MonitorFunc())
		}
	}
	return func() {
		cancelFunc()
		wg.Wait()
	}
}

func runMonitor(ctx context.Context, wg *sync.WaitGroup, period time.Duration, onTicker MonitorFunc) {
	go func() {
		wg.Add(1)
		defer wg.Done()
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				onTicker(ctx)
			}
		}
	}()
	return
}
