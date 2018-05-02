package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

func NewCGQuitGetCtrl() Result {
	const HTML = `
<html>
<head><title>Quit cgminer</title></head>
<body>
<form action="/cgminer/quit" method="POST">
<input type="submit" name="quit" value="Quit cgminer" style="
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
			Path:    "/cgminer/quit",
			Methods: []string{http.MethodGet},
			Handler: tool.HtmlHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(HTML))
			}),
		},
	}
}

func NewCGQuitPostCtrl(mgr monitor.Manager, client miner.Client) Result {
	return Result{
		Controller: &Controller{
			Path:    "/cgminer/quit",
			Methods: []string{http.MethodPost},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mgr.Stop()
				defer mgr.Start()

				bamStat := BAMStatus{Status: "OK"}
				httpStat := http.StatusOK

				err := client.Quit()
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

var CgminerQuitModule = fx.Provide(
	NewCGQuitGetCtrl,
	NewCGQuitPostCtrl,
)
