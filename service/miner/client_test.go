package miner

import (
	"testing"

	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"github.com/blockassets/bam_agent/tool"
)

func TestNewClient(t *testing.T) {
	config := cgminer.NewConfig(tool.NewCmdLine())
	client := NewClient(cgminer.NewConfigPort(config))
	err := client.Quit()
	if err == nil {
		t.Fatal("expected an error when calling quit!")
	}
}
