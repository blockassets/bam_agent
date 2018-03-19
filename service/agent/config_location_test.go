package agent

import (
	"testing"
)

func TestLocationHelper_Get(t *testing.T) {
	mc := NewMockConfig()
	cfg := NewConfigLocation(mc)

	data := cfg.Get()
	if data.Position != 1 {
		t.Fatalf("expected position to be 1, got %v", data.Position)
	}
}

func TestLocationHelper_Update(t *testing.T) {
	const newPosition = 5
	mc := NewMockConfig()
	cfg := NewConfigLocation(mc)
	data := cfg.Get()
	if data.Position != 1 {
		t.Fatalf("expected position to be 1, got %v", data.Position)
	}

	updateCfg := LocationConfig{
		Position: newPosition,
		}


	err := cfg.Update(updateCfg)
	if err != nil {
		t.Fatal(err)
	}

	if !mc.CalledSave {
		t.Fatalf("expected called save to be run, got %v", mc.CalledSave)
	}

	data = cfg.Get()
	if data.Position !=  updateCfg.Position {
		t.Fatalf("expected Position to be: %v, got: %v", newPosition, data.Position)
	}

	path := mc.Original().Path("location.position").Data()
	//json reads it as float64
	newPos := float64(newPosition)
	if path != newPos {
		t.Fatalf("expected position to be %v, got %v", newPosition, path)
	}
}
