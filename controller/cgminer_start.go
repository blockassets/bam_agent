package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/service"
	"github.com/json-iterator/go"
)

// Implements Builder interface
type CGStartCtrl struct{}

func (ctrl CGStartCtrl) build(cfg *Config) *Controller {
	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/cgminer/start",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl CGStartCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// we dont need to protect from reboot while doing this by stopping the monitors
			service.StartMiner()

			resp, _ := jsoniter.Marshal(BAMStatus{"OK", nil})
			w.Write(resp)
		})
}
