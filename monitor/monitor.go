package monitor

import (
	"context"
	"sync"
	"time"
)

type OnTick func() TickerFunc
type TickerFunc func(ctx context.Context)
type Stop func()

type Monitor interface {
	Start(ctx context.Context, wg *sync.WaitGroup, tickerFunc TickerFunc)
	IsEnabled() bool
	GetPeriod() time.Duration
	NewTickerFunc() TickerFunc
}

type Data struct {
	Enabled bool
	Period  time.Duration
	OnTick  OnTick
}

func (data *Data) Start(ctx context.Context, wg *sync.WaitGroup, tickerFunc TickerFunc) {
	go func() {
		wg.Add(1)
		defer wg.Done()
		ticker := time.NewTicker(data.Period)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				tickerFunc(ctx)
			}
		}
	}()
}

func (data *Data) IsEnabled() bool {
	return data.Enabled
}

func (data *Data) GetPeriod() time.Duration {
	return data.Period
}

func (data *Data) NewTickerFunc() TickerFunc {
	return data.OnTick()
}
