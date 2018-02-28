package monitor

import (
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service"
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

func (sr *testStatRetriever) GetLoad() (service.LoadAvgs, error) {
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

	return service.ParseLoad(data)
}

var countSomething int

func doSomething() { countSomething++ }

func TestCheckLoad(t *testing.T) {

	sr := &testStatRetriever{}
	sr.dataset = LevelNotEnough
	countSomething = 0
	tooHigh, err := checkLoad(sr, 5.0, doSomething)
	if err == nil {
		t.Errorf("t1.1: Expected error!")
	}
	if countSomething != 0 {
		t.Errorf("t1.1: Expected 0 in countSomething")
	}
	sr.dataset = LevelBelowFive
	countSomething = 0
	tooHigh, err = checkLoad(sr, 5.0, doSomething)
	if tooHigh {
		t.Errorf("t1.2: Expected low, got high!")
	}
	if countSomething != 0 {
		t.Errorf("t1.2: Expected 0 in countSomething")
	}
	sr.dataset = LevelExactlyFive
	countSomething = 0
	tooHigh, err = checkLoad(sr, 5.0, doSomething)
	if tooHigh {
		t.Errorf("t1.3: Expected low, got high!")
	}
	if countSomething != 0 {
		t.Errorf("t1.3: Expected 0 in countSomething")
	}
	sr.dataset = LevelAboveFive
	countSomething = 0
	tooHigh, err = checkLoad(sr, 5.0, doSomething)
	if !tooHigh {
		t.Errorf("t1.4: Expected high, got low!")
	}
	if countSomething != 1 {
		t.Errorf("t1.4: Expected 1 in countSomething")
	}
	sr.dataset = LevelMalformed
	countSomething = 0
	tooHigh, err = checkLoad(sr, 5.0, doSomething)
	if err == nil {
		t.Errorf("t1.5: Expected error!")
	}
	if countSomething != 0 {
		t.Errorf("t1.5: Expected 0 in countSomething")
	}
}

func TestLoadMonitors(t *testing.T) {
	testOnHighLoadCounter := 0

	sr := &testStatRetriever{}
	sr.dataset = LevelAboveFive
	cfg := MonitorConfig{}
	cfg.Load = LoadConfig{Enabled: true, PeriodSecs: 1, HighLoadMark: 5.0}

	lm := newLoadMonitor(sr, func() { testOnHighLoadCounter += 1 })

	err := lm.Start(&cfg)
	if err != nil {
		t.Errorf("t2.1 Expected start to suceed. Returned %+v", err)
	}
	if lm.getRunning() != true {
		t.Errorf("t2.2 Expected lm.isRunning to be true")
	}
	// give it time for one call
	time.Sleep(time.Duration(1500) * time.Millisecond)
	lm.Stop()
	if testOnHighLoadCounter != 1 {
		t.Errorf("t2.3 Expected 1 onHighMarks, got %d", testOnHighLoadCounter)
	}
	mark := testOnHighLoadCounter
	time.Sleep(time.Duration(1500) * time.Millisecond)
	if mark != testOnHighLoadCounter {
		t.Errorf("t2.4 Expected OnHighLoad to stop: mark == %d, counter == %d", mark, testOnHighLoadCounter)
	}
	err = lm.Start(&cfg)
	if err != nil {
		t.Errorf("t2.5 Expected 2nd start to suceed. Returned %+v", err)
	}
	err = lm.Start(&cfg)
	if err == nil {
		t.Errorf("t2.6 Expected 3rd start to fail")
	}
	lm.Stop()
	if lm.getRunning() {
		t.Errorf("t2.7 Expected to be not running")
	}
	cfg.Load.Enabled = false
	err = lm.Start(&cfg)
	if err != nil {
		t.Errorf("t2.8 Expected 4th start to succeed")
	}
	if !lm.getRunning() {
		t.Errorf("t2.9 Expected to be running")
	}
	lm.Stop()
	if lm.getRunning() {
		t.Errorf("t2.10 Expected to be not running")
	}

}
