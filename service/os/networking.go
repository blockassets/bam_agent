package os

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"go.uber.org/fx"
)

const (
	defaultNetworkInterfacesFile = "/etc/network/interfaces"
	staticIpInterfacesFmt        = `auto lo
auto eth0
iface eth0 inet static
address %s
netmask %s
gateway %s`
	dhcpInterfacesFmt = `auto lo
auto eth0
iface lo inet loopback
iface eth0 inet dhcp`
)

type Networking interface {
	SetDHCP() error
	SetStatic(address string, netmask string, gateway string) error
	FormatStatic(address string, netmask string, gateway string) string
}

type NetworkingData struct {
	File string
}

func (n *NetworkingData) SetDHCP() error {
	return ioutil.WriteFile(n.File, []byte(dhcpInterfacesFmt), 0644)
}

func (n *NetworkingData) SetStatic(address string, netmask string, gateway string) error {
	if (net.ParseIP(address) == nil) || (net.ParseIP(netmask) == nil) || (net.ParseIP(gateway) == nil) {
		return errors.New("invalid IP address")
	}
	return ioutil.WriteFile(n.File, []byte(n.FormatStatic(address, netmask, gateway)), 0644)
}

func (n *NetworkingData) FormatStatic(address string, netmask string, gateway string) string {
	return fmt.Sprintf(staticIpInterfacesFmt, address, netmask, gateway)
}

func NewNetworking() Networking {
	return &NetworkingData{File: defaultNetworkInterfacesFile}
}

var NetworkingModule = fx.Provide(NewNetworking)
