package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func NewConfigIPCtrl(mgr monitor.Manager, networking os.Networking, cfgNet miner.ConfigNetwork) Result {
	return Result{
		Controller: &Controller{
			Path:    "/config/ip",
			Methods: []string{http.MethodPut},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bamStat := BAMStatus{"OK", nil}
				httpStat := http.StatusOK

				// Declare things ahead of time to make the boolean logic below easier. grrrlang.
				var err error
				var data []byte

				data, err = ioutil.ReadAll(r.Body)
				if err == nil {
					mgr.Stop()
					defer mgr.Start()

					var netData *miner.NetworkData
					netData, err = cfgNet.Parse(data)

					if err == nil {
						err = cfgNet.Save(netData)
						if err == nil {
							err = networking.SetStatic(netData.IPAddress, netData.Netmask, netData.Gateway)
						}
					}
				}

				if err != nil {
					httpStat = http.StatusInternalServerError
					bamStat = BAMStatus{"Error", err}
				}

				w.WriteHeader(httpStat)
				resp, _ := jsoniter.Marshal(bamStat)
				w.Write(resp)
			}),
		},
	}
}
