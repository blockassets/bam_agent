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
	ctrl.Handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.Status != "200 OK" {
		t.Fatalf("expected OK and got %s", resp.Status)
	}

	bodyStr := string(body)
	if ! strings.Contains(bodyStr, "# HELP") {
		t.Fatalf("expected a # HELP and got %s", bodyStr)
	}
}
