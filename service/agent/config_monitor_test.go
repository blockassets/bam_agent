package agent

import (
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
)

func TestMonitorHelper_Get(t *testing.T) {
	mc := NewMockConfig()
	cfg := NewConfigMonitor(mc)

	data := cfg.Get()
	if data.Reboot.Enabled != false {
		t.Fatalf("expected reboot false, got %v", data.Reboot.Enabled)
	}
}

func TestMonitorHelper_Update(t *testing.T) {
	mc := NewMockConfig()
	cfg := NewConfigMonitor(mc)

	if cfg.Get().Reboot.Period.Duration < time.Duration(72)*time.Hour {
		t.Fatalf("expected > 72, got %v", cfg.Get().Reboot.Period.Duration)
	}

	updateCfg := MonitorConfig{
		Reboot: MonitorRebootConfig{
			Period: tool.RandomDuration{
				Duration: time.Duration(100) * time.Hour,
			},
		},
	}

	err := cfg.Update(updateCfg)
	if err != nil {
		t.Fatal(err)
	}

	if ! mc.CalledSave {
		t.Fatalf("expected called save to be run, got %v", mc.CalledSave)
	}

	data := cfg.Get()
	if data.Reboot.Period.Duration < time.Duration(100) * time.Hour {
		t.Fatalf("expected duration > 100, got: %v", data.Reboot.Period.Duration)
	}

	path := mc.Original().Path("monitor.reboot.period").Data().(string)
	if path != "100h0m0s" {
		t.Fatalf("expected period to be 100h0m0s, got %s", path)
	}
}
