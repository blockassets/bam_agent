package tool

import (
	"strings"
	"testing"
)

const (
	jsonA = `{"foo": {"jon": 1}}`
	jsonB = `{"foo": {"jon": 2}, "bar": {"hello": 1}}`
)

func TestMerge(t *testing.T) {
	result, err := Merge([]byte(jsonA), []byte(jsonB))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(result), "hello") {
		t.Fatalf("expected hello in the string, got %v", result)
	}

	if !strings.Contains(string(result), "jon\":1") {
		t.Fatalf("expected jon:1 in the string, got %v", string(result))
	}
}
