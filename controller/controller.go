package controller

import (
	"net/http"

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

type Builder interface {
	build() *Controller
	makeHandler() http.HandlerFunc
}

func makeJsonHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8") // normal header
		handler.ServeHTTP(w, r)
	}
}

func Init(e *echo.Echo) {
	ctrls := []*Controller{RebootCtrl{}.build(), CGQuitCtrl{}.build()}

	for _, ctrl := range ctrls {
		e.Match(ctrl.Methods, ctrl.Path, echo.WrapHandler(ctrl.Handler))
	}
}
