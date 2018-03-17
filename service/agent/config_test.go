package agent

import (
	"testing"

	"github.com/blockassets/bam_agent/tool"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig(tool.NewCmdLine())

	if ! cfg.Loaded().Monitor.HighLoad.Enabled {
		t.Fatalf("expected highLoad to be enabled")
	}
}
