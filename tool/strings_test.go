package tool

import "testing"

func TestTrimToNil(t *testing.T) {
	str := TrimToNil("")
	if str != nil {
		t.Fatalf("expected nil, got '%v'", *str)
	}

	str = TrimToNil("   ")
	if str != nil {
		t.Fatalf("expected nil, got '%v'", *str)
	}
}
