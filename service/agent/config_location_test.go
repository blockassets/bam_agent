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
	const newRack = "rack1"
	const newPosition = 5
	const invalidPosition = 0
	const defaultShelf = 1
	const invalidShelf = 0

	mc := NewMockConfig()
	cfg := NewConfigLocation(mc)
	data := cfg.Get()
	if data.Position != 1 {
		t.Fatalf("expected position to be 1, got %v", data.Position)
	}

	updateCfg := LocationConfig{
		Rack:     newRack,
		Position: newPosition,
		Shelf:    defaultShelf,
	}

	err := cfg.Update(updateCfg)
	if err != nil {
		t.Fatal(err)
	}

	if !mc.CalledSave {
		t.Fatalf("expected called save to be run, got %v", mc.CalledSave)
	}

	data = cfg.Get()
	if data.Position != updateCfg.Position {
		t.Fatalf("expected Position to be: %v, got: %v", newPosition, data.Position)
	}

	if data.Rack != updateCfg.Rack {
		t.Fatalf("expected Rack to be: %v, got: %v", newRack, data.Rack)
	}

	path := mc.Original().Path("location.position").Data()
	//json reads it as float64
	newPos := float64(newPosition)
	if path != newPos {
		t.Fatalf("expected position to be %v, got %v", newPosition, path)
	}

	invalidPositionCfg := LocationConfig{
		Position: invalidPosition,
		Shelf:    defaultShelf,
	}
	err = cfg.Update(invalidPositionCfg)
	if err == nil {
		t.Fatal("Expected an error on an invlalid position")
	}
	invalidShelfCfg := LocationConfig{
		Position: newPosition,
		Shelf:    invalidShelf,
	}
	err = cfg.Update(invalidShelfCfg)
	if err == nil {
		t.Fatal("Expected an error on an invlalid shelf")
	}

}
