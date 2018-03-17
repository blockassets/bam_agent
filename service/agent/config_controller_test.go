package agent

import (
	"testing"
	"time"
)

func TestControllerHelper_Get(t *testing.T) {
	mc := NewMockConfig()
	cfg := NewConfigController(mc)

	data := cfg.Get()
	if data.Reboot.Delay != time.Duration(5)*time.Second {
		t.Fatalf("expected reboot false, got %v", data.Reboot.Delay)
	}
}

func TestControllerHelper_Update(t *testing.T) {
	mc := NewMockConfig()
	cfg := NewConfigController(mc)

	if cfg.Get().Reboot.Delay != time.Duration(5)*time.Second {
		t.Fatalf("expected == 5s, got %v", cfg.Get().Reboot.Delay)
	}

	updateCfg := ControllerConfig{
		Reboot: ControllerRebootConfig{
			Delay: time.Duration(100) * time.Hour,
		},
	}

	err := cfg.Update(updateCfg)
	if err != nil {
		t.Fatal(err)
	}

	if !mc.CalledSave {
		t.Fatalf("expected called save to be run, got %v", mc.CalledSave)
	}

	data := cfg.Get()
	if data.Reboot.Delay != time.Duration(100)*time.Hour {
		t.Fatalf("expected duration > 100, got: %v", data.Reboot.Delay)
	}

	path := mc.Original().Path("controller.reboot.delay").Data().(string)
	if path != "100h0m0s" {
		t.Fatalf("expected period to be 100h0m0s, got %s", path)
	}
}
