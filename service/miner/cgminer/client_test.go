package cgminer

import (
	"testing"
)

func TestCGMinerClient_Quit(t *testing.T) {
	cgClient := NewMockCgClient()
	client := NewClientWrapper(cgClient)

	err := client.Quit()
	if err != nil {
		t.Fatal("Did not expected an error when calling quit!")
	}
	if !cgClient.quitCalled {
		t.Fatal("Expected quitCalled to be true!")
	}
}

func TestCGMinerClient_GetAccepted(t *testing.T) {
	cgClient := NewMockCgClient()
	client := NewClientWrapper(cgClient)

	accepted, err := client.GetAccepted()
	if err != nil {
		t.Fatal("Did not expected an error when calling GetAccepted!")
	}
	if accepted != testAccepted {
		t.Fatalf("Expected accepted to be equal to testAccepted, got %v", accepted)
	}
}

func TestCGMinerClient_GetTemp(t *testing.T) {
	cgClient := NewMockCgClient()
	client := NewClientWrapper(cgClient)

	temp, err := client.GetTemp()
	if err != nil {
		t.Fatal("Did not expected an error when calling GetTemp!")
	}
	if temp != testTemp {
		t.Fatalf("Expected accepted to be equal to testTemp, got %v", temp)
	}
}
