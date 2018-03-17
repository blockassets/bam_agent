package os

import (
	"net"
	"testing"
)

func TestNetInfoData_GetMacAddress(t *testing.T) {

	ni := &NetInfoData{
		netInterfaceFunc: func() ([]net.Interface, error) {
			return NewNetInterfaces(), nil
		},
	}

	addr := ni.GetMacAddress()
	if *addr != "66:6f:6f:6f" {
		t.Fatalf("expected 66:6f:6f:6f, got %s", *addr)
	}
}
