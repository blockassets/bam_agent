package tool

import (
	"testing"
)

func TestNewCmdLine(t *testing.T) {
	cmdLine := NewCmdLine()
	if cmdLine.MinerConfigPath != "/usr/app/conf.default" {
		t.Fatalf("expected /usr/app/conf.default got %s", cmdLine.MinerConfigPath)
	}

	if cmdLine.Port != ":1111" {
		t.Fatalf("expected :1111 got %s", cmdLine.Port)
	}

	if cmdLine.AgentConfigPath != "/etc/bam_agent.json" {
		t.Fatalf("expected /etc/bam_agent.json got %s", cmdLine.AgentConfigPath)
	}
}
