package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/service"
)

// Implements Controller interface
type RebootCtrl struct{}

func (c RebootCtrl) build() *Controller {
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
			cmds := service.Command{}
			go cmds.Reboot()
		})
}
