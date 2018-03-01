package main

import (
	"os"
	"testing"
)

func TestDefaultConfFilePath(t *testing.T) {
	_, err := os.Open("conf/bam_agent.json")
	if os.IsNotExist(err) {
		t.Fatal(err)
	}
}
