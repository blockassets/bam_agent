package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/agent"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func NewPutLocationCtrl(mgr monitor.Manager, location agent.ConfigLocation) Result {
	return Result{
		Controller: &Controller{
			Path:    "/config/location",
			Methods: []string{http.MethodPut},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bamStat := BAMStatus{"OK", nil, ""}
				httpStat := http.StatusOK

				// Declare things ahead of time to make the boolean logic below easier. grrrlang.
				var err error
				var data []byte
				var inLocation *agent.LocationConfig

				data, err = ioutil.ReadAll(r.Body)
				if err == nil {
					mgr.Stop()
					defer mgr.Start()

					inLocation, err = location.Parse(data)
					if err == nil {
						err = location.Update(*inLocation)
					}
				}

				if err != nil {
					httpStat = http.StatusInternalServerError
					bamStat = BAMStatus{"Error", err, ""}
				}

				w.WriteHeader(httpStat)
				resp, _ := jsoniter.Marshal(bamStat)
				w.Write(resp)
			}),
		},
	}
}
