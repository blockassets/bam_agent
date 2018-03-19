package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func NewCGStartCtrl(mgr monitor.Manager, miner os.Miner) Result {
	return Result{
		Controller: &Controller{
			Path:    "/cgminer/start",
			Methods: []string{http.MethodGet},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bamStat := BAMStatus{Status: "OK"}
				httpStat := http.StatusOK

				mgr.Stop()
				defer mgr.Start()

				err := miner.Start()
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
