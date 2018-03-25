package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewNodeExporterCtrl(t *testing.T) {
	result := NewNodeExporterCtrl()
	ctrl := result.Controller

	if ctrl.Path != "/metrics/node_exporter" {
		t.Fatalf("expected /metrics/node_exporter, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodGet {
		t.Fatalf("expected method get, got %s", ctrl.Methods[0])
	}

	req := httptest.NewRequest("GET", "/doesnotmatter", nil)
	w := httptest.NewRecorder()
	nodeExporterHandler([]string{"cpu"}).ServeHTTP(w, req)

	resp := w.Result()
	if resp.Status != "200 OK" {
		t.Fatalf("expected OK and got %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "# TYPE node_cpu_seconds_total counter") {
		t.Fatalf("expected '# TYPE node_cpu_seconds_total counter', got %s", bodyStr)
	}
}
