package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/service/os"
	"github.com/json-iterator/go"
)

func TestNewRebootCtrl_Post(t *testing.T) {
	cfg := RebootConfig{Delay: time.Duration(50) * time.Millisecond}

	reboot := os.NewMockReboot()
	result := NewRebootPostCtrl(cfg, &reboot)
	ctrl := result.Controller

	if ctrl.Path != "/reboot" {
		t.Fatalf("expected /reboot, got %s", ctrl.Path)
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

	time.Sleep(cfg.Delay * 2)

	if reboot.Counter != 1 {
		t.Fatalf("expected counter 1 and got %d", reboot.Counter)
	}
}
