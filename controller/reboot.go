package controller

import (
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/service"
	"github.com/json-iterator/go"
)

const (
	DELAY_BEFORE_REBOOT = time.Duration(5) * time.Second
)

// Implements Builder interface
type RebootCtrl struct{}

func (ctrl RebootCtrl) build(cfg *Config) *Controller {
	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/reboot",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl RebootCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			resp, _ := jsoniter.Marshal(BAMStatus{"OK", nil})
			w.Write(resp)
			// leave enough time for http server to respond to caller
			time.AfterFunc(DELAY_BEFORE_REBOOT, service.Reboot)
		})
}
