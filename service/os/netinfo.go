package os

import (
	"net"

	"go.uber.org/fx"
)

// The miners always have a single physical network.
// Return the first physical MAC address enumerated as the primary MAC for the device

type NetInfo interface {
	GetNetInterfaces() ([]net.Interface, error)
	GetMacAddress() *string
}

type NetInfoData struct {
	netInterfaceFunc func() ([]net.Interface, error)
}

func (ni *NetInfoData) GetNetInterfaces() ([]net.Interface, error) {
	return ni.netInterfaceFunc()
}

func (ni *NetInfoData) GetMacAddress() *string {
	interfaces, err := ni.GetNetInterfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && len(i.HardwareAddr) != 0 {
				mac := i.HardwareAddr.String()
				return &mac
			}
		}
	}
	return nil
}

func NewNetInfo() NetInfo {
	return &NetInfoData{
		netInterfaceFunc: net.Interfaces,
	}
}

var NetInfoModule = fx.Provide(NewNetInfo)

// For testing only
func NewNetInterfaces() []net.Interface {
	noFace := net.Interface{Flags: net.FlagUp}
	oneFace := net.Interface{HardwareAddr: []byte("fooo"), Flags: net.FlagUp}
	twoFace := net.Interface{HardwareAddr: []byte("bar")}
	threeFace := net.Interface{HardwareAddr: []byte("ack"), Flags: net.FlagUp}
	expected := []net.Interface{noFace, oneFace, twoFace, threeFace}

	return expected
}
