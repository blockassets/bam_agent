package monitor

import (
	"testing"
	"time"
)

type testStatRetriever struct {
	dataset int
}

const (
	LevelNotEnough = iota
	LevelBelowFive
	LevelExactlyFive
	LevelAboveFive
	LevelMalformed
)

func (sr *testStatRetriever) getLoad() (LoadAvgs, error) {
	var data string
	switch sr.dataset {
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

	return parseLoad(data)

}

func doNothing() {}

func TestcheckLoad(t *testing.T) {
	sr := &testStatRetriever{}
	sr.dataset = LevelNotEnough
	tooHigh, err := checkLoad(sr, 5.0, doNothing)
	if err == nil {
		t.Errorf("Expected error!")
	}
	sr.dataset = LevelBelowFive
	tooHigh, err = checkLoad(sr, 5.0, doNothing)
	if tooHigh {
		t.Errorf("Expected low, got high!")
	}
	sr.dataset = LevelExactlyFive
	tooHigh, err = checkLoad(sr, 5.0, doNothing)
	if tooHigh {
		t.Errorf("Expected low, got high!")
	}
	sr.dataset = LevelAboveFive
	tooHigh, err = checkLoad(sr, 5.0, doNothing)
	if !tooHigh {
		t.Errorf("Expected high, got low!")
	}
	sr.dataset = LevelMalformed
	tooHigh, err = checkLoad(sr, 5.0, doNothing)
	if err == nil {
		t.Errorf("Expected error!")
	}
}

func TestLoadMonitors(t *testing.T) {
	testOnHighLoadCounter := 0

	sr := &testStatRetriever{}
	sr.dataset = LevelBelowFive
	cfg := MonitorConfig{}
	cfg.Load = LoadConfig{Enabled: true, PeriodSecs: 1, HighLoadMark: 5.0}

	lm := newLoadMonitor(sr, func() { testOnHighLoadCounter += 1 })

	err := lm.Start(&cfg)
	if err != nil {
		t.Errorf("Expected start to suceed. Returned %+v", err)
	}
	if lm.getRunning() != true {
		t.Errorf("Expected lm.isRunning to be true")
	}
	time.Sleep(time.Duration(4500) * time.Millisecond)
	if testOnHighLoadCounter != 0 {
		t.Errorf("Expected 0 onHighMarks, got %d", testOnHighLoadCounter)
	}
	sr.dataset = LevelAboveFive
	time.Sleep(time.Duration(4500) * time.Millisecond)
	lm.Stop()
	mark := testOnHighLoadCounter
	time.Sleep(time.Duration(2000) * time.Millisecond)
	if mark != testOnHighLoadCounter {
		t.Errorf("Expected OnHighLoad to stop: mark == %d, counter == %d", mark, testOnHighLoadCounter)
	}

	err = lm.Start(&cfg)
	if err != nil {
		t.Errorf("Expected 2nd start to suceed. Returned %+v", err)
	}
	err = lm.Start(&cfg)
	if err == nil {
		t.Errorf("Expected 3rd start to fail")
	}

}
