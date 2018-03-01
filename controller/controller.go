package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/cgminer_client"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
)

var (
	json = jsoniter.ConfigDefault
)

type BAMStatus struct {
	Status string
	Error  error
}

type Controller struct {
	Methods []string
	Path    string
	Handler http.HandlerFunc
}

type Config struct {
	Version string
	Client  *cgminer_client.Client
	MonitorManager *monitor.Manager
}

type Builder interface {
	build(cfg *Config) *Controller
	makeHandler() http.HandlerFunc
}

func makeJsonHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Prevent caching of any of the requests so that we can use GET for things like /reboot
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Allow any origin to hit the miner. Potential danger here, but if someone has access
		// to your miner network, then you've got bigger problems anyway since they have a default
		// root password (and some machines even have telnet enabled).
		w.Header().Set("Access-Control-Allow-Origin", "*")

		handler.ServeHTTP(w, r)
	}
}

func Init(e *echo.Echo, cfg *Config) {
	ctrls := []*Controller{
		RebootCtrl{}.build(cfg),
		CGQuitCtrl{}.build(cfg),
		PutPoolsCtrl{}.build(cfg),
		StatusCtrl{}.build(cfg),
	}

	for _, ctrl := range ctrls {
		e.Match(ctrl.Methods, ctrl.Path, echo.WrapHandler(ctrl.Handler))
	}
}
