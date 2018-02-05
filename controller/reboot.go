package controller

import (
	"net/http"
	"os/exec"
)

// Implements Controller interface
type RebootCtrl struct {}

func (c RebootCtrl) build() *Controller {
	return &Controller{
		Methods: []string{http.MethodGet},
		Path: "/reboot",
		Handler: c.makeHandler(),
	}
}

func (c RebootCtrl) makeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exec.Command("reboot", "-f")
	}
}
