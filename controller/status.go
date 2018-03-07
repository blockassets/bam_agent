package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/blockassets/bam_agent/service"
	"github.com/json-iterator/go"
)

// Implements Builder interface
type StatusCtrl struct {
	version string
}

type Status struct {
	Agent      string        `json:"agent"`
	Miner      string        `json:"miner"`
	Uptime     time.Duration `json:"uptime"`
	MACAddress *string       `json:"mac"`
}

func (ctrl StatusCtrl) build(cfg *Config) *Controller {
	ctrl.version = cfg.Version

	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/status",
		Handler: ctrl.makeHandler(),
	}
}

func (ctrl StatusCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {

			uptime, _ := service.GetUptime()
			ni := service.NewNetInfo()
			ni.GetMacAddress()

			status := Status{
				Agent:  strings.TrimSpace(ctrl.version),
				Miner:  strings.TrimSpace(service.ReadVersionFile()),
				Uptime: uptime,
				MACAddress: ni.GetMacAddress(),
			}

			w.WriteHeader(http.StatusOK)
			resp, _ := jsoniter.Marshal(status)
			w.Write(resp)
		})
}
