package controller

import (
	"log"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/cgminer_client"
	"github.com/json-iterator/go"
)

// Implements Builder interface
type CGQuitCtrl struct {
	client         *cgminer_client.Client
	monitorManager *monitor.Manager
}

func (ctrl CGQuitCtrl) build(cfg *Config) *Controller {
	ctrl.client = cfg.Client
	ctrl.monitorManager = cfg.MonitorManager

	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/cgminer/quit",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl CGQuitCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			log.Printf("CGMiner Quit Requested")

			bamStat := BAMStatus{"OK", nil}
			httpStat := http.StatusOK

			ctrl.monitorManager.StopMonitors()

			err := ctrl.client.Quit()
			if err != nil {
				httpStat = http.StatusBadGateway
				bamStat = BAMStatus{"Error", err}
			}

			ctrl.monitorManager.StartMonitors()

			w.WriteHeader(httpStat)
			resp, _ := jsoniter.Marshal(bamStat)
			w.Write(resp)
		})
}
