package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
)

const (
	interfacesFmt = `auto lo
auto eth0
iface eth0 inet static
address %s
netmask %s
gateway %s`
	defaultNetworkInterfacesFile = "/etc/network/interfaces"
)

// Miners have ipv4 addresses and configuration
func SetStaticIP(address string, netmask string, gateway string) error {
	err := writeInterfacesFile(defaultNetworkInterfacesFile, address, netmask, gateway)
	return err
}

func writeInterfacesFile(filename string, address string, netmask string, gateway string) error {
	if (net.ParseIP(address) == nil) || (net.ParseIP(netmask) == nil) || (net.ParseIP(gateway) == nil) {
		return errors.New("invalid IP address")
	}
	out := fmt.Sprintf(interfacesFmt, address, netmask, gateway)
	err := ioutil.WriteFile(filename, []byte(out), 0644)
	return err
}
