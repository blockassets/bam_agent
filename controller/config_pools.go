package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

func NewPutPoolsCtrl(mgr monitor.Manager, poolCfg miner.ConfigPools, client miner.Client) Result {
	return Result{
		Controller: &Controller{
			Path:    "/config/pools",
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

					var pools *miner.PoolAddresses
					pools, err = poolCfg.Parse(data)
					if err == nil {
						err = poolCfg.Save(pools)
						if err == nil {
							err = client.Quit()
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

func NewGetPoolsCtrl(poolCfg miner.ConfigPools) Result {
	return Result{
		Controller: &Controller{
			Path:    "/config/pools",
			Methods: []string{http.MethodGet},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var response interface{}
				var httpStat int

				pools, err := poolCfg.Get()
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
			}),
		},
	}
}

var ConfigPoolsModule = fx.Provide(
	NewPutPoolsCtrl,
	NewGetPoolsCtrl,
)
