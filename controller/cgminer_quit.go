package controller

import (
	"log"
	"net/http"

	"github.com/blockassets/bam_agent/service"
)

// Implements Controller interface
type CGQuitCtrl struct{}

func (c CGQuitCtrl) build() *Controller {
	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/cgminer/quit",
		Handler: c.makeHandler(),
	}
}

func (c CGQuitCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CGMiner Quit Requested")

			bamStat := BAMStatus{"OK", nil}
			httpStat := http.StatusOK

			err := service.CgmQuit()
			if err != nil {
				httpStat = http.StatusBadGateway
				bamStat = BAMStatus{"Error", err}
			}
			w.WriteHeader(httpStat)
			resp, _ := json.Marshal(bamStat)
			w.Write(resp)
		})
}
