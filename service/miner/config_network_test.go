package miner

import (
	"testing"
)

const (
	netConfig = `{"ip": "10.10.0.11", "mask": "255.255.252.0", "gateway": "10.10.0.1", "dns": "8.8.8.8"}`
)

func TestNetworkHelper_Parse(t *testing.T) {
	nw := &NetworkHelper{}
	res, err := nw.Parse([]byte(netConfig))
	if err != nil {
		t.Fatal(err)
	}
	if res.IPAddress != "10.10.0.11" {
		t.Fatalf("expected ip 10.10.0.11, got %s", res.IPAddress)
	}

	if res.Netmask != "255.255.252.0" {
		t.Fatalf("expected ip 255.255.252.0, got %s", res.Netmask)
	}

	if res.Netmask != "255.255.252.0" {
		t.Fatalf("expected ip 10.10.0.11, got %s", res.IPAddress)
	}

	if res.Dns != "8.8.8.8" {
		t.Fatalf("expected dns 8.8.8.8, got %s", res.Dns)
	}
}

func TestNetworkHelper_Save(t *testing.T) {
	cfg := NewMockConfig(DefaultConfigFile)
	nw := &NetworkHelper{Config: &cfg}
	res, err := nw.Parse([]byte(netConfig))
	if err != nil {
		t.Fatal(err)
	}
	err = nw.Save(res)
	if err != nil {
		t.Fatal(err)
	}

	res, err = nw.Parse([]byte(cfg.Saved))
	if err != nil {
		t.Fatal(err)
	}

	if res.Dns != "8.8.8.8" {
		t.Fatalf("expected dns 8.8.8.8, got %s", res.Dns)
	}
}

func TestNetworkHelper_Get(t *testing.T) {
	cfg := NewMockConfig(netConfig)
	nw := &NetworkHelper{Config: &cfg}
	res, err := nw.Get()
	if err != nil {
		t.Fatal(err)
	}

	if res.Dns != "8.8.8.8" {
		t.Fatalf("expected dns 8.8.8.8, got %s", res.Dns)
	}
}
