package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service"
	"github.com/json-iterator/go"
)

// Implements Builder interface
type PutDhcpCtrl struct {
	monitorManager *monitor.Manager
}

func (ctrl PutDhcpCtrl) build(cfg *Config) *Controller {
	ctrl.monitorManager = cfg.MonitorManager

	return &Controller{
		Methods: []string{http.MethodPut},
		Path:    "/config/dhcp",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl PutDhcpCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			bamStat := BAMStatus{"OK", nil}
			httpStat := http.StatusOK

			ctrl.monitorManager.Stop()

			err := service.UpdateDHCPNetConfig()
			if err != nil {
				httpStat = http.StatusBadGateway
				bamStat = BAMStatus{"Error", err}
			}

			ctrl.monitorManager.Start()

			w.WriteHeader(httpStat)
			resp, _ := jsoniter.Marshal(bamStat)
			w.Write(resp)
		})
}
