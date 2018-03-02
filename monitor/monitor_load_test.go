package monitor

import (
	"sync"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type testStatRetriever struct {
	dataSet int
}

const (
	LevelNotEnough = iota
	LevelBelowFive
	LevelExactlyFive
	LevelAboveFive
	LevelMalformed
)

func (sr *testStatRetriever) GetLoad() (service.LoadAvgs, error) {
	var data string
	switch sr.dataSet {
	case LevelNotEnough:
		data = "0.0 0.0"
	case LevelBelowFive:
		data = "0.0 4.999 0.0 1234 1234"
	case LevelExactlyFive:
		data = "0.0 5.0 0.0 1234 1234"
	case LevelAboveFive:
		data = "0.0 5.1 0.0 1234 1234"
	case LevelMalformed:
		data = "a b c d emnf,masfd"
	}

	return service.ParseLoad(data)
}

func NewTestStatRetriever(dataSet int) service.StatRetriever {
	return &testStatRetriever{
		dataSet: dataSet,
	}
}

var countSomething int

func doSomething() { countSomething++ }

func TestCheckLoad(t *testing.T) {
	expectErrors := []service.StatRetriever{
		NewTestStatRetriever(LevelNotEnough),
		NewTestStatRetriever(LevelMalformed),
	}

	for _, sr := range expectErrors {
		countSomething = 0
		_, err := checkLoad(sr, 5.0, doSomething)
		if err == nil {
			t.Errorf("Expected error, got nil!")
		}
		if countSomething != 0 {
			t.Errorf("Expected 0 in countSomething, got %d", countSomething)
		}
	}

	expectLow := []service.StatRetriever{
		NewTestStatRetriever(LevelBelowFive),
		NewTestStatRetriever(LevelExactlyFive),
	}

	for _, sr := range expectLow {
		countSomething = 0
		tooHigh, _ := checkLoad(sr, 5.0, doSomething)
		if tooHigh {
			t.Errorf("Expected low, got high!")
		}
		if countSomething != 0 {
			t.Errorf("Expected 0 in countSomething, got %d", countSomething)
		}
	}

	expectHigh := []service.StatRetriever{
		NewTestStatRetriever(LevelAboveFive),
	}

	for _, sr := range expectHigh {
		countSomething = 0
		tooHigh, _ := checkLoad(sr, 5.0, doSomething)
		if !tooHigh {
			t.Errorf("Expected high, got low!")
		}
		if countSomething != 1 {
			t.Errorf("Expected 1 in countSomething, got %d", countSomething)
		}
	}
}

func TestLoadMonitor_Start(t *testing.T) {
	count := 0

	config := &HighLoadConfig{Enabled: true, PeriodInSeconds: 1, HighLoadMark: 5.0}
	sr := NewTestStatRetriever(LevelAboveFive)
	tickerPeriod := time.Duration(50) * time.Millisecond

	onHighLoad := func() { count += 1 }

	context := &Context{quit: make(chan bool), waitGroup: &sync.WaitGroup{}}

	monitor := newLoadMonitor(context, config, &tickerPeriod, sr, onHighLoad)
	err := monitor.Start()
	if err != nil {
		t.Error(err)
	}

	// Sleep to ensure the timer runs once
	time.Sleep(tickerPeriod * 2)

	// Test that stop cleans up the WaitGroup
	monitor.Stop()
	context.waitGroup.Wait()

	if count == 0 {
		t.Errorf("Expected >=1 count, got %d", count)
	}
}
