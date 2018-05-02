package controller

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"github.com/json-iterator/go"
)

const (
	testPoolData = `{"pool1": "1", "pool2": "2", "pool3": "3"}`
)

func TestNewConfigPoolsCtrl(t *testing.T) {

	client := miner.NewMockMinerClient(-1)
	cfg := cgminer.NewMockConfig(cgminer.DefaultConfigFile)
	ph := cgminer.NewPoolHelper(&cfg)
	mgr := monitor.NewMockManager()

	result := NewPutPoolsCtrl(&mgr, ph, &client)
	ctrl := result.Controller

	if ctrl.Path != "/config/pools" {
		t.Fatalf("expected /config/pools, got %s", ctrl.Path)
	}

	if len(ctrl.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(ctrl.Methods))
	}

	if ctrl.Methods[0] != http.MethodPut {
		t.Fatalf("expected method put, got %s", ctrl.Methods[0])
	}

	put := strings.NewReader(testPoolData)
	req := httptest.NewRequest("PUT", "/doesnotmatter", put)
	w := httptest.NewRecorder()
	ctrl.Handler.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	status := &BAMStatus{}
	err := jsoniter.Unmarshal(body, status)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	if resp.Header.Get("Expires") != "0" {
		t.Fatalf("Expires: 0, got %s", resp.Header.Get("Expires"))
	}

	if status.Status != "OK" {
		t.Fatalf("expected OK and got %s", status.Status)
	}

	res, err := ph.Get()
	if err != nil {
		t.Fatal(err)
	}

	if res.Pool1 != "1" {
		t.Fatalf("expected pool1 = 1, got %s", res.Pool1)
	}

	if client.CalledRestart != true {
		t.Fatalf("expected client.CalledRestart and got %v", client.CalledRestart)
	}

	if !mgr.CalledStop && !mgr.CalledStart {
		t.Fatalf("expected manager stop/start, got start: %v, stop: %v", mgr.CalledStart, mgr.CalledStop)
	}
}

func TestNewGetPoolsCtrl(t *testing.T) {
	cfg := cgminer.NewMockConfig(testPoolData)
	ph := cgminer.NewPoolHelper(&cfg)

	result := NewGetPoolsCtrl(ph)
	ctrl := result.Controller

	if ctrl.Path != "/config/pools" {
		t.Fatalf("expected /config/pools, got %s", ctrl.Path)
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

	pools, err := ph.Parse(body)
	if err != nil {
		t.Fatal(err)
	}

	if pools.Pool1 != "1" {
		t.Fatalf("expected pool1 = 1, got %s", pools.Pool1)
	}
}
