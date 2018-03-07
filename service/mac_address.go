package service

import (
	"bytes"
	"net"
)

// The miners we have always have  a single physical network.
// So we take the first physical MAC address enumerated as the primary MAC for the device
//

const (
	nullMACAddress = "00:00:00:00:00:00"
)

type NetInfo struct {
	ifi *[]net.Interface
}

func NewNetInfo() *NetInfo {
	ifi, err := net.Interfaces()
	if err != nil {
		return &NetInfo{nil}
	} else {
		return &NetInfo{&ifi}
	}
}

func (ni *NetInfo) GetPrimaryMacAddress() string {
	if ni.ifi != nil {
		for _, i := range *ni.ifi {
			if (i.Flags&net.FlagUp != 0) && (bytes.Compare(i.HardwareAddr, nil) != 0) {
				return i.HardwareAddr.String()
			}
		}
	}
	return nullMACAddress
}
