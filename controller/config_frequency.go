package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func NewConfigFrequencyCtrl(mgr monitor.Manager, cfgFreq miner.ConfigFrequency, client miner.Client) Result {
	return Result{
		Controller: &Controller{
			Path:    "/config/frequency",
			Methods: []string{http.MethodPut},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bamStat := BAMStatus{Status: "OK"}
				httpStat := http.StatusOK

				// Declare things ahead of time to make the boolean logic below easier. grrrlang.
				var err error
				var data []byte

				data, err = ioutil.ReadAll(r.Body)
				if err == nil {
					mgr.Stop()
					defer mgr.Start()

					var freq *miner.FrequencyData
					freq, err = cfgFreq.Parse(data)
					if err == nil {
						err = cfgFreq.Save(freq.Frequency)
						if err == nil {
							err = client.Quit()
						}
					}
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
