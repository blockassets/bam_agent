package os

import (
	"testing"
	"time"
)

func TestParseUptime(t *testing.T) {
	timeStr := "1304750.82 1938868.54"

	result := parseUptime(timeStr)

	if result.Duration != time.Duration(1304750) {
		t.Fatalf("Expected 1304750, but got %s", result)
	}

	result = parseUptime("123")
	if result.Error == nil {
		t.Fatal(result.Error)
	}

	result = parseUptime("123 123")
	if result.Error == nil {
		t.Fatal(result.Error)
	}
}
