package service

import (
	"net"
	"testing"
)

func makeNetErrNetInfo() *NetInfo {
	return &NetInfo{nil}
}

func makeNoPhysicalNetInfo() *NetInfo {
	ifs := &[]net.Interface{{0, 1500, "test_1", net.HardwareAddr{}, net.FlagUp & net.FlagLoopback}}
	return &NetInfo{ifs}
}

func TestNetInfo_GetPrimaryMacAddress(t *testing.T) {
	// Real Network
	ni := NewNetInfo()
	mac := ni.GetMacAddress()
	if mac == nil {
		t.Errorf("Expected a valid MAC address")
	}

	ni = makeNetErrNetInfo()
	mac = ni.GetMacAddress()
	if mac != nil {
		t.Errorf("Expected nil, got: %v", mac)
	}

	ni = makeNoPhysicalNetInfo()
	mac = ni.GetMacAddress()
	if mac != nil {
		t.Errorf("Expected nil, got: %v", mac)
	}
}
