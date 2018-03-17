package os

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewNetworking(t *testing.T) {
	address := "10.10.0.1"
	netmask := "255.255.240.0"
	gateway := "10.10.0.2"

	file, err := ioutil.TempFile("/tmp", "network-test")
	defer file.Close()
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	networking := &NetworkingData{File: file.Name()}
	networking.SetDHCP()

	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != dhcpInterfacesFmt {
		t.Fatalf("expected dhcpInterfacesFmt, got: %s", string(data))
	}

	err = networking.SetStatic(address, netmask, gateway)
	if err != nil {
		t.Fatal(err)
	}

	data, err = ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != networking.FormatStatic(address, netmask, gateway) {
		t.Fatalf("expected staticIpInterfacesFmt, got: %s", string(data))
	}
}
