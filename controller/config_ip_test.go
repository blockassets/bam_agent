package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner/cgminer"
	sos "github.com/blockassets/bam_agent/service/os"
	"github.com/json-iterator/go"
)

func TestNewConfigIPCtrl(t *testing.T) {
	file, err := ioutil.TempFile("", "network-ip")
	defer file.Close()
	defer os.Remove(file.Name())

	if err != nil {
		t.Fatal(err)
	}

	networking := &sos.NetworkingData{File: file.Name()}
	cfg := cgminer.NewMockConfig("")
	cfgNet := cgminer.NewConfigNetwork(&cfg)
	mgr := monitor.NewMockManager()

	result := NewConfigIPCtrl(&mgr, networking, cfgNet)
	ctrl := result.Controller

	if ctrl.Path != "/config/ip" {
		t.Fatalf("expected /config/ip, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPut {
		t.Fatalf("expected method get, got %s", ctrl.Methods[0])
	}

	putData := `{"ip": "10.10.0.11", "mask": "255.255.252.0", "gateway": "10.10.0.1", "dns": "8.8.8.8"}`
	req := httptest.NewRequest("PUT", "/doesnotmatter", strings.NewReader(putData))
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
