package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/cgminer_client"
	"github.com/json-iterator/go"
)

// command to update the ip address for
// pools
// PUT
//{
//"pool1":"",
//"pool2":"",
//"pool3":""
//}
// eg { "pool1":"111.2.3.4", "pool2":"112.3.4.5", "pool3":"113.4.5.6"}
// and we update the conf.default file on the miner

// Implements Builder interface
type PutPoolsCtrl struct {
	client         *cgminer_client.Client
	monitorManager *monitor.Manager
}

func (ctrl PutPoolsCtrl) build(cfg *Config) *Controller {
	ctrl.client = cfg.Client
	ctrl.monitorManager = cfg.MonitorManager

	return &Controller{
		Methods: []string{http.MethodPut},
		Path:    "/config/pools",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl PutPoolsCtrl) makeHandler() http.HandlerFunc {
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

				err = service.UpdatePools(data)
				if err != nil {
					httpStat = http.StatusBadGateway
					bamStat = BAMStatus{"Error", err}
				}

				err = ctrl.client.Quit()
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

type GetPoolsCtrl struct{}

func (ctrl GetPoolsCtrl) build(cfg *Config) *Controller {

	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/config/pools",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl GetPoolsCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {
			var response interface{}
			var httpStat int

			pools, err := service.GetPools()
			if err != nil {
				httpStat = http.StatusInternalServerError
				response = BAMStatus{"Error", err}
			} else {
				httpStat = http.StatusOK
				response = pools
			}

			w.WriteHeader(httpStat)
			resp, _ := jsoniter.Marshal(response)
			w.Write(resp)
		})
}
