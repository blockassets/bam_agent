package agent

import (
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfig(tool.NewCmdLine())
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Controller.Reboot.Delay != time.Duration(5)*time.Second {
		t.Fatalf("expected 5s and got %s", cfg.Controller.Reboot)
	}
}
