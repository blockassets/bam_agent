package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/agent"
	"github.com/json-iterator/go"
)

const (
	testLocationData = `{
    "facility": "",
    "row": "",
    "shelf": 5,
    "position": 5
}`
)

func TestNewConfigLocationCtrl(t *testing.T) {

	mc := agent.NewMockConfig()
	mgr := monitor.NewMockManager()
	cfg := agent.NewConfigLocation(mc)

	result := NewPutLocationCtrl(&mgr, cfg)
	ctrl := result.Controller

	if ctrl.Path != "/config/location" {
		t.Fatalf("expected /config/location, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPut {
		t.Fatalf("expected method put, got %s", ctrl.Methods[0])
	}

	put := strings.NewReader(testLocationData)
	req := httptest.NewRequest("PUT", "/doesnotmatter", put)
	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	status := &BAMStatus{}
	err := jsoniter.Unmarshal(body, status)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	if resp.Header.Get("Expires") != "0" {
		t.Fatalf("Expires: 0, got %s", resp.Header.Get("Expires"))
	}

	if status.Status != "OK" {
		t.Fatalf("expected OK and got %s", status.Status)
	}

	loc := cfg.Get()
	if loc.Position != 5 {
		t.Fatalf("expected Position = 5, got %v", loc.Position)
	}
}
