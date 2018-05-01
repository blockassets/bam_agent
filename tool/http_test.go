package tool

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonHandlerFunc_ServeHTTP(t *testing.T) {
	req := httptest.NewRequest("GET", "/doesnotmatter", nil)
	w := httptest.NewRecorder()

	var calledFunc = false

	fun := JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledFunc = true
	})
	fun.ServeHTTP(w, req)

	resp := w.Result()

	if len(resp.Header) != 5 {
		t.Fatalf("expected 4 headers, got %v", len(resp.Header))
	}

	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("got wrong content-type header")
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

	if ! calledFunc {
		t.Fatalf("didn't call the function")
	}
}

func TestJsonHandlerFunc_ServeHTTP_NoPurpose(t *testing.T) {
	req := httptest.NewRequest("GET", "/doesnotmatter", nil)
	req.Header.Set("X-Purpose", "Preview")

	w := httptest.NewRecorder()

	var calledFunc = false

	fun := JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledFunc = true
	})
	fun.ServeHTTP(w, req)

	resp := w.Result()

	if len(resp.Header) != 1 {
		t.Fatalf("expected 1 headers, got %v", len(resp.Header))
	}

	if resp.Header.Get("Content-Type") != "text/plain" {
		t.Fatalf("got wrong content-type header")
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("got wrong status header")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if len(body) != len("No preview allowed") {
		t.Fatalf("got wrong body")
	}

	if calledFunc {
		t.Fatalf("shouldn't call the function")
	}
}
