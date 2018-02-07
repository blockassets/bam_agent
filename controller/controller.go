package controller

import (
	"net/http"
	"github.com/labstack/echo"
)

type BAMStatus struct {
	Status string
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



func Init(e *echo.Echo) {
	// TODO: Make this more automated once there are more controllers
	ctrl := RebootCtrl{}.build()
	e.Match(ctrl.Methods, ctrl.Path, echo.WrapHandler(ctrl.Handler))
}
