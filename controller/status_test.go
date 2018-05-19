package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/agent"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/json-iterator/go"
)

func TestNewStatusCtrl(t *testing.T) {
	uptimeFunc := func() os.UptimeResult {
		return os.UptimeResult{
			Duration: time.Duration(42) * time.Second,
		}
	}

	netInfo := os.NewMockNetInfo()
	mc := agent.NewMockConfig()
	cfg := agent.NewConfigLocation(mc)

	result := NewStatusCtrl(agent.Version{V: "1"}, miner.Version{V: "2"}, uptimeFunc, &netInfo, cfg)
	ctrl := result.Controller

	if ctrl.Path != "/status" {
		t.Fatalf("expected /status, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodGet {
		t.Fatalf("expected method get, got %s", ctrl.Methods[0])
	}

	req := httptest.NewRequest("GET", "/doesnotmatter", nil)
	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	status := &StatusResponse{}
	err := jsoniter.Unmarshal(body, status)
	if err != nil {
		t.Fatal(err)
	}

	if *status.Agent != "1" {
		t.Fatalf("expected 1 and got %d", status.Agent)
	}

	if *status.Miner != "2" {
		t.Fatalf("expected 2 and got %d", status.Miner)
	}

	if status.Uptime != time.Duration(42)*time.Second {
		t.Fatalf("expected 42s and got %s", status.Uptime)
	}

	if status.Mac == nil {
		t.Fatalf("expected not nil for Mac, got %v", status.Mac)
	}

	if status.Location.Position != 1 {
		t.Fatalf("expected 1 for location position, got %v", status.Location.Position)
	}

	if !status.Date.Before(time.Now()) {
		t.Fatalf("expected Date to be before now(), got %v", status.Date)
	}
}
