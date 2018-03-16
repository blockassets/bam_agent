package tool

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonHandlerFunc_ServeHTTP(t *testing.T) {
	req := httptest.NewRequest("GET", "/doesnotmatter", nil)
	w := httptest.NewRecorder()

	fun := JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fun.ServeHTTP(w, req)

	resp := w.Result()

	if len(resp.Header) != 5 {
		t.Fatalf("expected 4 headers, got %v", len(resp.Header))
	}

	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("got wrong cache-control header")
	}

	if resp.Header.Get("Cache-Control") != "no-cache, no-store, must-revalidate" {
		t.Fatalf("got wrong cache-control header")
	}

	if resp.Header.Get("Pragma") != "no-cache" {
		t.Fatalf("got wrong pragma header")
	}

	if resp.Header.Get("Expires") != "0" {
		t.Fatalf("got wrong expires header")
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("got wrong ACAO header")
	}
}
