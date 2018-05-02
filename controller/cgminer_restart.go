package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)


func NewCGRestartGetCtrl(mgr monitor.Manager, client miner.Client) Result {
	const HTML = `
<html>
<head><title>Restart cgminer</title></head>
<body>
<form action="/cgminer/restart" method="POST">
<input type="submit" name="restart" value="Restart cgminer" style="
	width: 50%;
	margin: 30px;
	font-size: 25px;
	padding: 30px;
	border-radius: 12px;
	background-color: #f00;
	color: #fff;" />
</form>
</body>
</html>
`

	return Result{
		Controller: &Controller{
			Path:    "/cgminer/restart",
			Methods: []string{http.MethodGet},
			Handler: tool.HtmlHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(HTML))
			}),
		},
	}
}

func NewCGRestartPostCtrl(mgr monitor.Manager, client miner.Client) Result {
	return Result{
		Controller: &Controller{
			Path:    "/cgminer/restart",
			Methods: []string{http.MethodPost},
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

var CgminerRestartModule = fx.Provide(
	NewCGRestartGetCtrl,
	NewCGRestartPostCtrl,
)
