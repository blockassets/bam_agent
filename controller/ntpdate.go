package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

func NewNtpdateGetCtrl() Result {
	const HTML = `
<html>
<head><title>Ntpdate</title></head>
<body>
<form action="/ntpdate" method="POST">
<input type="submit" name="ntpdate" value="Ntpdate" style="
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
			Path:    "/ntpdate",
			Methods: []string{http.MethodGet},
			Handler: tool.HtmlHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(HTML))
			}),
		},
	}
}

func NewNtpdatePostCtrl(ntpdate os.Ntpdate) Result {
	return Result{
		Controller: &Controller{
			Path:    "/ntpdate",
			Methods: []string{http.MethodPost},
			Handler: tool.JsonHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var resp []byte

				before := time.Now()
				err := ntpdate.Ntpdate()
				after := time.Now()

				message := fmt.Sprintf("Before: %s, After: %s", before, after)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					resp, _ = jsoniter.Marshal(BAMStatus{Status: "Error", Error: err, Message: message})
				} else {
					w.WriteHeader(http.StatusOK)
					resp, _ = jsoniter.Marshal(BAMStatus{Status: "OK", Message: message})
				}

				w.Write(resp)
			}),
		},
	}
}

var NtpdateModule = fx.Provide(
	NewNtpdateGetCtrl,
	NewNtpdatePostCtrl,
)
