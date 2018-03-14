package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service"
	"github.com/json-iterator/go"
)

// Implements Builder interface
type PutFrequencyCtrl struct {
	monitorManager *monitor.Manager
}

func (ctrl PutFrequencyCtrl) build(cfg *Config) *Controller {
	ctrl.monitorManager = cfg.MonitorManager

	return &Controller{
		Methods: []string{http.MethodPut},
		Path:    "/config/frequency",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl PutFrequencyCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			bamStat := BAMStatus{"OK", nil}
			httpStat := http.StatusOK

			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				httpStat = http.StatusInternalServerError
				bamStat = BAMStatus{"Error", err}
			} else {
				ctrl.monitorManager.Stop()

				err = service.UpdateFrequency(data)
				if err != nil {
					httpStat = http.StatusBadGateway
					bamStat = BAMStatus{"Error", err}
				}

				ctrl.monitorManager.Start()
			}

			w.WriteHeader(httpStat)
			resp, _ := jsoniter.Marshal(bamStat)
			w.Write(resp)
		})
}
