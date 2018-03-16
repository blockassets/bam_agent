package miner

import (
	"testing"
)

const (
	portConfig = `{"api-port": "4099"}`
)

func TestPortHelper_Get(t *testing.T) {
	cfg := NewMockConfig(portConfig)
	ph := &PortHelper{Config: &cfg}
	if ph.Get() != 4099 {
		t.Fatalf("expected 4099, got %v", ph.Get())
	}
}
