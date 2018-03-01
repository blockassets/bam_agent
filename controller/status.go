package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/blockassets/bam_agent/service"
)

// Implements Builder interface
type StatusCtrl struct {
	version string
}

type Status struct {
	Agent  string        `json:"agent"`
	Miner  string        `json:"miner"`
	Uptime time.Duration `json:"uptimeInSeconds"`
}

func (c StatusCtrl) build(cfg *Config) *Controller {
	c.version = cfg.Version

	return &Controller{
		Methods: []string{http.MethodGet},
		Path:    "/status",
		Handler: c.makeHandler(),
	}
}

func (c StatusCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {

			uptime, _ := service.GetUptime()

			status := Status{
				Agent:  strings.TrimSpace(c.version),
				Miner:  strings.TrimSpace(service.ReadVersionFile()),
				Uptime: uptime,
			}

			w.WriteHeader(http.StatusOK)
			resp, _ := json.Marshal(status)
			w.Write(resp)
		})
}
