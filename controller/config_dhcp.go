package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func NewConfigDHCPCtrl(mgr monitor.Manager, networking os.Networking, cfgNet cgminer.ConfigNetwork) Result {
	return Result{
		Controller: &Controller{
			Path:    "/config/dhcp",
			Methods: []string{http.MethodPut},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bamStat := BAMStatus{Status: "OK"}
				httpStat := http.StatusOK

				mgr.Stop()
				defer mgr.Start()

				err := cfgNet.Save(&cgminer.NetworkData{})
				if err == nil {
					err = networking.SetDHCP()
				}

				if err != nil {
					httpStat = http.StatusInternalServerError
					bamStat = BAMStatus{Status: "Error", Error: err}
				}

				w.WriteHeader(httpStat)
				resp, _ := jsoniter.Marshal(bamStat)
				w.Write(resp)
			}),
		},
	}
}
