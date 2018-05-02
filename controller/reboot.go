package controller

import (
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

type RebootConfig struct {
	Delay time.Duration
}

func NewRebootGetCtrl() Result {
	const HTML = `
<html>
<head><title>Reboot</title></head>
<body>
<form action="/reboot" method="POST">
<input type="submit" name="reboot" value="Reboot" style="
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
			Path:    "/reboot",
			Methods: []string{http.MethodGet},
			Handler: tool.HtmlHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(HTML))
			}),
		},
	}
}

func NewRebootPostCtrl(cfg RebootConfig, reboot os.Reboot) Result {
	return Result{
		Controller: &Controller{
			Path:    "/reboot",
			Methods: []string{http.MethodPost},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)

				resp, _ := jsoniter.Marshal(BAMStatus{Status: "OK"})
				w.Write(resp)
				// leave enough time for http server to respond to caller
				time.AfterFunc(cfg.Delay, func() { reboot.Reboot() })
			}),
		},
	}
}

var RebootModule = fx.Provide(
	NewRebootGetCtrl,
	NewRebootPostCtrl,
)
