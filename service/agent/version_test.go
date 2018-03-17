package agent

import "testing"

func TestNewVersion(t *testing.T) {
	v := NewVersion()
	if len(v.V) > 0 {
		t.Fatalf("expected no value for version, got %s", v.V)
	}

	version = "foo"

	v = NewVersion()
	if v.V != "foo" {
		t.Fatalf("expected foo and got %s", v.V)
	}
}
