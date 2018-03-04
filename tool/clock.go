package tool

import "time"

/*
	A simple interface wrapper around Timer/Ticker so that we can pass a Clock
	around instead of having to duplicate logic.
*/
type Clock interface {
	C() <-chan time.Time
	Stop()
}

type Ticker struct {
	ticker *time.Ticker
}

type Timer struct {
	timer *time.Timer
}

func (ticker Ticker) C() <-chan time.Time {
	return ticker.ticker.C
}

func (ticker Ticker) Stop() {
	ticker.ticker.Stop()
}

func (timer Timer) C() <-chan time.Time {
	return timer.timer.C
}

func (timer Timer) Stop() {
	timer.timer.Stop()
}

func NewTicker(period time.Duration) Clock {
	return &Ticker{ticker: time.NewTicker(period)}
}

func NewTimer(period time.Duration) Clock {
	return &Timer{timer: time.NewTimer(period)}
}
