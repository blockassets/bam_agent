package cgminer

import (
	"testing"
)

const (
	test1Pool = `{"pool1": "1", "pool2": "2", "pool3": "3"}`
)

func TestPoolHelper_Parse(t *testing.T) {
	p := &PoolHelper{}
	pools, err := p.Parse([]byte(test1Pool))
	if err != nil {
		t.Fatal(err)
	}

	if pools.Pool1 != "1" {
		t.Fatalf("expected 1 and got %s", pools.Pool1)
	}
}

func TestPoolHelper_Save(t *testing.T) {
	cfg := NewMockConfig(DefaultConfigFile)
	ph := &PoolHelper{Config: &cfg}

	poolAddress, err := ph.Parse([]byte(test1Pool))
	if err != nil {
		t.Fatal(err)
	}

	err = ph.Save(poolAddress)
	if err != nil {
		t.Fatal(err)
	}

	poolAddressSaved, err := ph.Parse([]byte(cfg.Saved))

	if poolAddressSaved.Pool1 != "1" {
		t.Fatalf("expected 1, got %s", poolAddressSaved.Pool1)
	}
}

func TestPoolHelper_Get(t *testing.T) {
	cfg := NewMockConfig(test1Pool)
	ph := &PoolHelper{Config: &cfg}
	addr, err := ph.Get()

	if err != nil {
		t.Fatal(err)
	}

	if addr.Pool1 != "1" {
		t.Fatalf("expected 1, got %s", addr.Pool1)
	}
}
