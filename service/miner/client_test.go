package miner

import (
	"testing"

	"github.com/blockassets/bam_agent/tool"
)

func TestNewClient(t *testing.T) {
	config := NewConfig(tool.NewCmdLine())
	client := NewClient(NewConfigPort(config))
	err := client.Quit()
	if err == nil {
		t.Fatal("expected an error when calling quit!")
	}
}
