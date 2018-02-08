package monitor

import (
	"testing"
)

type testStatRetriever struct {
	dataset int
}

func (sr testStatRetriever) getLoad() (LoadAvgs, error) {
	var data string
	switch sr.dataset {
	case 1:
		data = "0.0 0.0" // not enough
	case 2:
		data = "0.0 4.999 0.0 1234 1234"
	case 3:
		data = "0.0 5.0 0.0 1234 1234"
	case 4:
		data = "0.0 5.1 0.0 1234 1234"
	case 5:
		data = "a b c d emnf,masfd"
	}

	return parseLoad(data)

}

func TestMonitorLoad(t *testing.T) {
	sr := testStatRetriever{}
	sr.dataset = 1
	too_high, err := checkLoadAvg(sr)
	if err == nil {
		t.Errorf("Expected error!")
	}
	sr.dataset = 2
	too_high, err = checkLoadAvg(sr)
	if too_high {
		t.Errorf("Expected low, got high!")
	}
	sr.dataset = 3
	too_high, err = checkLoadAvg(sr)
	if too_high {
		t.Errorf("Expected low, got high!")
	}
	sr.dataset = 4
	too_high, err = checkLoadAvg(sr)
	if !too_high {
		t.Errorf("Expected high, got low!")
	}
	sr.dataset = 5
	too_high, err = checkLoadAvg(sr)
	if err == nil {
		t.Errorf("Expected error!")
	}
}
