package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

func NewCGStartGetCtrl() Result {
	const HTML = `
<html>
<head><title>Start cgminer</title></head>
<body>
<form action="/cgminer/start" method="POST">
<input type="submit" name="start" value="Start cgminer" style="
	width: 50%;
	margin: 30px;
	font-size: 25px;
	padding: 30px;
	border-radius: 12px;
	background-color: green;
	color: #fff;" />
</form>
</body>
</html>
`

	return Result{
		Controller: &Controller{
			Path:    "/cgminer/start",
			Methods: []string{http.MethodGet},
			Handler: tool.HtmlHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(HTML))
			}),
		},
	}
}

func NewCGStartPostCtrl(mgr monitor.Manager, miner os.Miner) Result {
	return Result{
		Controller: &Controller{
			Path:    "/cgminer/start",
			Methods: []string{http.MethodPost},
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

var CgminerStartModule = fx.Provide(
	NewCGStartGetCtrl,
	NewCGStartPostCtrl,
)
