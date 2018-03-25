package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/miner/cgminer"
)

/*
	Impossible to test fully without doing a mock server
*/
func TestNewCgminerExporterCtrl(t *testing.T) {
	mockCgClient := cgminer.NewMockCgClient()

	result := NewCgMinerExporterCtrl(mockCgClient, miner.NewVersion(nil))
	ctrl := result.Controller

	if ctrl.Path != "/metrics/cgminer_exporter" {
		t.Fatalf("expected /metrics/cgminer_exporter, got %s", ctrl.Path)
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
	if resp.Status != "200 OK" {
		t.Fatalf("expected OK and got %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "cgminer_summary_work_utility") {
		t.Fatalf("expected cgminer_summary_work_utility, got %s", bodyStr)
	}
}
