package main

import (
	"os"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestDefaultConfFilePath(t *testing.T) {
	_, err := os.Open("conf/bam_agent.json")
	if os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestLoadAgentConfig(t *testing.T) {
	// Normally done for us in main.go
	tool.RegisterTimeDuration()
	tool.RegisterRandomDuration()

	outputFile := "/tmp/bam_agent.json"

	config, err := LoadAgentConfig(outputFile)
	if err != nil {
		t.Fatal(err)
	}

	// Test the Random Duration loading of the file
	duration, _ := time.ParseDuration("24h")
	if config.Monitor.CGMQuit.Period.Duration <= duration {
		t.Fatal(config.Monitor.CGMQuit.Period.Duration)
	}

	duration2, _ := time.ParseDuration("72h")
	if config.Monitor.Reboot.Period.Duration <= duration2 {
		t.Fatal(config.Monitor.Reboot.Period.Duration)
	}

	// Cleanup the temporary test file
	err = os.Remove(outputFile)
	if err != nil {
		t.Fatal(err)
	}
}
