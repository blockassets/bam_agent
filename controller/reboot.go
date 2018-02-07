package controller

import (
	"time"
	"net/http"
	"os/exec"
	"encoding/json"
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
	return func(w http.ResponseWriter, r *http.Request) {
		bamStat := BAMStatus{"OK"}
		resp, _ := json.Marshal(bamStat)
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // normal header
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		go reboot()
	}
}

func reboot() {
	time.Sleep(5 * time.Second)
	exec.Command("/sbin/reboot", "-f").Run()

}