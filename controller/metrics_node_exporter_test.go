package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Sadly, travis doesn't have /proc/stat so NE fails there. Thus, we can't do full testing.
*/
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
	nodeExporterHandler().ServeHTTP(w, req)

	resp := w.Result()
	if resp.Status != "200 OK" {
		t.Fatalf("expected OK and got %s", resp.Status)
	}
}
