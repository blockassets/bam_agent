package service

import (
	"testing"
	"time"
)

func TestParseUptime(t *testing.T) {
	timeStr := "1304750.82 1938868.54"

	result, err := parseUptime(timeStr)
	if err != nil {
		t.Error(err)
	}

	if result != time.Duration(1304750) {
		t.Errorf("Expected 1304750, but got %s", result)
	}

	result, err = parseUptime("123")
	if err == nil {
		t.Error(err)
	}

	result, err = parseUptime("123 123")
	if err == nil {
		t.Error(err)
	}
}
