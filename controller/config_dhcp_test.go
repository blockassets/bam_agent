package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	sos "github.com/blockassets/bam_agent/service/os"
	"github.com/json-iterator/go"
)

func TestNewConfigDHCPCtrl(t *testing.T) {
	file, err := ioutil.TempFile("", "network-dhcp")
	defer file.Close()
	defer os.Remove(file.Name())

	if err != nil {
		t.Fatal(err)
	}

	networking := &sos.NetworkingData{File: file.Name()}
	cfg := miner.NewMockConfig("")
	cfgNet := miner.NewConfigNetwork(&cfg)
	mgr := monitor.NewMockManager()

	result := NewConfigDHCPCtrl(&mgr, networking, cfgNet)
	ctrl := result.Controller

	if ctrl.Path != "/config/dhcp" {
		t.Fatalf("expected /config/dhcp, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPut {
		t.Fatalf("expected method get, got %s", ctrl.Methods[0])
	}

	req := httptest.NewRequest("PUT", "/doesnotmatter", nil)
	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	status := &BAMStatus{}
	err = jsoniter.Unmarshal(body, status)
	if err != nil {
		t.Fatal(err)
	}

	if status.Status != "OK" {
		t.Fatalf("expected OK and got %s", status.Status)
	}

	data, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Don't really care about contents of file because lower tests do that
	if len(data) == 0 {
		t.Fatalf("expected data, got %s", string(data))
	}

	if len(cfg.Saved) == 0 {
		t.Fatalf("expected cfg.Saved to have data, got %v", len(cfg.Saved))
	}

	if !mgr.CalledStop && !mgr.CalledStart {
		t.Fatalf("expected manager stop/start, got start: %v, stop: %v", mgr.CalledStart, mgr.CalledStop)
	}
}
