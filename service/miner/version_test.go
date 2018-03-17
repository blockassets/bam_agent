package miner

import (
	"testing"

	"github.com/blockassets/bam_agent/service/agent"
	"github.com/blockassets/bam_agent/tool"
)

func TestNewVersion(t *testing.T) {
	cfg := agent.NewConfig(tool.CmdLine{})
	version := NewVersion(cfg)

	if len(version.V) != 0 {
		t.Fatalf("expected no version since we aren't running on a miner")
	}
}
