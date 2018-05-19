package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blockassets/bam_agent/service/os"
	"github.com/json-iterator/go"
)

func TestNewNtpdateCtrl_Post(t *testing.T) {
	ntpdate := os.NewMockNtpdate()
	result := NewNtpdatePostCtrl(&ntpdate)
	ctrl := result.Controller

	if ctrl.Path != "/ntpdate" {
		t.Fatalf("expected /ntpdate, got %s", ctrl.Path)
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

	if ntpdate.Counter != 1 {
		t.Fatalf("expected counter 1 and got %d", ntpdate.Counter)
	}
}
