package controller

import (
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

type RebootConfig struct {
	Delay time.Duration
}

func NewRebootCtrl(cfg RebootConfig, reboot os.Reboot) Result {
	return Result{
		Controller: &Controller{
			Path:    "/reboot",
			Methods: []string{http.MethodGet},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				resp, _ := jsoniter.Marshal(BAMStatus{"OK", nil})
				w.Write(resp)
				// leave enough time for http server to respond to caller
				time.AfterFunc(cfg.Delay, func() { reboot.Reboot() })
			}),
		},
	}
}
