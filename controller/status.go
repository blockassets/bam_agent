package controller

import (
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/service/agent"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

type StatusResponse struct {
	Agent    *string              `json:"agent"`
	Miner    *string              `json:"miner"`
	Uptime   time.Duration        `json:"uptime"`
	Date     time.Time            `json:"date"`
	Mac      *string              `json:"mac"`
	Location agent.LocationConfig `json:"location"`
}

func NewStatusCtrl(agentVersion agent.Version, minerVersion miner.Version, getUptimeResult os.UptimeResultFunc, netInfo os.NetInfo, location agent.ConfigLocation) Result {
	return Result{
		Controller: &Controller{
			Path:    "/status",
			Methods: []string{http.MethodGet},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				uptime := getUptimeResult()

				status := StatusResponse{
					Agent:    tool.TrimToNil(agentVersion.V),
					Miner:    tool.TrimToNil(minerVersion.V),
					Uptime:   uptime.Duration,
					Mac:      netInfo.GetMacAddress(),
					Location: location.Get(),
					Date:     time.Now(),
				}

				w.WriteHeader(http.StatusOK)
				resp, _ := jsoniter.Marshal(status)
				w.Write(resp)
			}),
		},
	}
}
