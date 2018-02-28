package controller

import (
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/service"
)

const (
	DELAY_BEFORE_REBOOT = time.Duration(5) * time.Second
)

// Implements Builder interface
type RebootCtrl struct{}

func (c RebootCtrl) build(cfg *Config) *Controller {
	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/reboot",
		Handler: c.makeHandler(),
	}
}

func (c RebootCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			resp, _ := json.Marshal(BAMStatus{"OK", nil})
			w.Write(resp)
			// leave enough time for http server to respond to caller
			time.AfterFunc(DELAY_BEFORE_REBOOT, service.Reboot)
		})
}
