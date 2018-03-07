package service

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

const (
	expectedResult = `auto lo
auto eth0
iface eth0 inet static
address 1.1.1.1
netmask 2.2.2.2
gateway 3.3.3.3`
)

func TestWriteInterfacesFile(t *testing.T) {

	outputFile := "/tmp/interfaces"
	err := writeInterfacesFile(outputFile, "1.1.1.1", "2.2.2.2", "3.3.3.3")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	out, err := ioutil.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	// Cleanup the temporary test file
	err = os.Remove(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(out, []byte(expectedResult)) != 0 {
		t.Errorf("unexpected output, got %v", out)
	}

	err = writeInterfacesFile(outputFile, "", "", "")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
