package controller

import (
	"log"
	"net/http"
	"os/exec"
	"time"
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

			go reboot()
		})
}

func reboot() {
	time.Sleep(5 * time.Second)
	log.Printf("Reboot Requested")
	exec.Command("/sbin/reboot", "-f").Run()
}
