package miner

import (
	"testing"
)

const (
	freqConfig = `{"frequency": "100"}`
)

func TestFrequencyHelper_Get(t *testing.T) {
	cfg := NewMockConfig(freqConfig)
	ph := &FrequencyHelper{Config: &cfg}
	if ph.Get() != 100 {
		t.Fatalf("expected 100, got %v", ph.Get())
	}
}

func TestFrequencyHelper_Save(t *testing.T) {
	cfg := NewMockConfig(DefaultConfigFile)
	nw := &FrequencyHelper{Config: &cfg}
	res, err := nw.Parse([]byte(freqConfig))
	if err != nil {
		t.Fatal(err)
	}
	err = nw.Save(res.Frequency)
	if err != nil {
		t.Fatal(err)
	}

	res, err = nw.Parse([]byte(cfg.Saved))
	if err != nil {
		t.Fatal(err)
	}

	if res.Frequency != 100 {
		t.Fatalf("expected 100, got %v", res.Frequency)
	}

}
