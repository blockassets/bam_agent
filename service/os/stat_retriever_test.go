package os

import "testing"

func TestLinuxStatRetriever_GetLoadData(t *testing.T) {
	counter := 0
	lsr := &LinuxStatRetriever{
		getProcData: func(loadPath string) ([]byte, error) {
			counter++
			return nil, nil
		},
	}

	lsr.GetLoadData()

	if counter == 0 {
		t.Fatalf("expected counter 1, got %v", counter)
	}
}
