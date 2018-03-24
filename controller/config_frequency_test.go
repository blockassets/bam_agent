package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"github.com/json-iterator/go"
)

func TestNewConfigFrequencyCtrl(t *testing.T) {
	cfg := cgminer.NewMockConfig("")
	cfgFreq := cgminer.NewConfigFrequency(&cfg)
	client := miner.NewMockMinerClient(-1)
	mgr := monitor.NewMockManager()

	result := NewConfigFrequencyCtrl(&mgr, cfgFreq, &client)
	ctrl := result.Controller

	if ctrl.Path != "/config/frequency" {
		t.Fatalf("expected /config/frequency, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPut {
		t.Fatalf("expected method get, got %s", ctrl.Methods[0])
	}

	putData := `{"frequency": "648"}`
	req := httptest.NewRequest("PUT", "/doesnotmatter", strings.NewReader(putData))
	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	status := &BAMStatus{}
	err := jsoniter.Unmarshal(body, status)
	if err != nil {
		t.Fatal(err)
	}

	if status.Status != "OK" {
		t.Fatalf("expected OK and got %s", status.Status)
	}

	if len(cfg.Saved) == 0 {
		t.Fatalf("expected cfg.Saved to have data, got %v", len(cfg.Saved))
	}

	if !client.CalledQuit {
		t.Fatalf("expected client.Quit(), got %v", client.CalledQuit)
	}

	if !mgr.CalledStop && !mgr.CalledStart {
		t.Fatalf("expected manager stop/start, got start: %v, stop: %v", mgr.CalledStart, mgr.CalledStop)
	}
}
