package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func NewCGRestartCtrl(mgr monitor.Manager, client miner.Client) Result {
	return Result{
		Controller: &Controller{
			Path:    "/cgminer/restart",
			Methods: []string{http.MethodGet},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mgr.Stop()
				defer mgr.Start()

				bamStat := BAMStatus{Status: "OK"}
				httpStat := http.StatusOK

				err := client.Restart()
				if err != nil {
					httpStat = http.StatusBadGateway
					bamStat = BAMStatus{Status: "Error", Error: err}
				}

				w.WriteHeader(httpStat)
				resp, _ := jsoniter.Marshal(bamStat)
				w.Write(resp)
			}),
		},
	}
}
