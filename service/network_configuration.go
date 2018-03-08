package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
)

const (

	defaultNetworkInterfacesFile = "/etc/network/interfaces"
	staticIpInterfacesFmt = `auto lo
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

// Miners have ipv4 addresses and configuration
func SetInterfaceToStaticIP(address string, netmask string, gateway string) error {
	return writeStaticIpInterfaces(defaultNetworkInterfacesFile, address, netmask, gateway)
}

func writeStaticIpInterfaces(filename string, address string, netmask string, gateway string) error {
	if (net.ParseIP(address) == nil) || (net.ParseIP(netmask) == nil) || (net.ParseIP(gateway) == nil) {
		return errors.New("invalid IP address")
	}
	out := fmt.Sprintf(staticIpInterfacesFmt, address, netmask, gateway)
	return ioutil.WriteFile(filename, []byte(out), 0644)

}

func SetInterfaceToDhcp() error {
	return writeDhcpInterfaces(defaultNetworkInterfacesFile)
}

func writeDhcpInterfaces(filename string) error {
	return ioutil.WriteFile(filename, []byte(dhcpInterfacesFmt), 0644)
}
