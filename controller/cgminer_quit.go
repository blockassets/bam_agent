package controller

import (
	"log"
	"net/http"

	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/cgminer_client"
)

// Implements Builder interface
type CGQuitCtrl struct {
	client *cgminer_client.Client
}

func (c CGQuitCtrl) build(cfg *Config) *Controller {
	c.client = cfg.Client
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

			err := service.CgmQuit(c.client)
			if err != nil {
				httpStat = http.StatusBadGateway
				bamStat = BAMStatus{"Error", err}
			}
			w.WriteHeader(httpStat)
			resp, _ := json.Marshal(bamStat)
			w.Write(resp)
		})
}
