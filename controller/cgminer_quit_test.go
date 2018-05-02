package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/json-iterator/go"
)

func TestNewCGQuitPostCtrl(t *testing.T) {
	client := miner.NewMockMinerClient(-1)
	mgr := monitor.NewMockManager()
	result := NewCGQuitPostCtrl(&mgr, &client)
	ctrl := result.Controller

	if ctrl.Path != "/cgminer/quit" {
		t.Fatalf("expected /cgminer/quit, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPost {
		t.Fatalf("expected method post, got %s", ctrl.Methods[0])
	}

	req := httptest.NewRequest("POST", "/doesnotmatter", nil)
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

	if client.CalledQuit != true {
		t.Fatalf("expected calledQuit and got %v", client.CalledQuit)
	}

	if !mgr.CalledStop && !mgr.CalledStart {
		t.Fatalf("expected manager stop/start, got start: %v, stop: %v", mgr.CalledStart, mgr.CalledStop)
	}
}
