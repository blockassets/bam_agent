package os

import (
	"testing"
)

func TestLinuxMemInfo_Get(t *testing.T) {
	counter := 0
	memInfo := &LinuxMemInfo{
		getData: func(path string) ([]byte, error) {
			counter++
			return nil, nil
		},
	}

	_, err := memInfo.Get()
	if err != nil {
		t.Fatalf("expected error, got %v", err)
	}

	if counter == 0 {
		t.Fatalf("expected counter 1, got %v", counter)
	}
}
