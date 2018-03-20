package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/service/agent"
	"github.com/labstack/echo"
	"go.uber.org/fx"
)

type BAMStatus struct {
	Status  string
	Error   error
	Message string
}

type Controller struct {
	Methods []string
	Path    string
	Handler http.Handler
}

type Result struct {
	fx.Out
	Controller *Controller `group:"controller"`
}

type Group struct {
	fx.In
	Controllers []*Controller `group:"controller"`
}

type Manager interface {
	Match()
}

type Data struct {
	Controllers []*Controller
	MatchFunc   func(methods []string, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) []*echo.Route
}

func (d *Data) Match() {
	for _, c := range d.Controllers {
		d.MatchFunc(c.Methods, c.Path, echo.WrapHandler(c.Handler))
	}
}

/*
	There is a little fx magic here. Group gets magically populated with a list of
	controllers because those are provided in the module declaration below and we
	use the fx 'group' functionality to make that happen.
*/
func NewManager(e *echo.Echo, g Group) Manager {
	data := &Data{
		MatchFunc:   e.Match,
		Controllers: g.Controllers,
	}
	data.Match()
	return data
}

var Module = fx.Options(

	// Copy the config to keep the dependencies clean
	fx.Provide(func(cfg agent.ControllerConfig) RebootConfig {
		return RebootConfig{
			Delay: cfg.Reboot.Delay,
		}
	}),

	// Separate since it provides the two controllers
	ConfigPoolsModule,

	fx.Provide(
		NewManager,

		NewCGQuitCtrl,
		NewCGStartCtrl,
		NewConfigDHCPCtrl,
		NewConfigFrequencyCtrl,
		NewConfigIPCtrl,
		NewGetPoolsCtrl,
		NewRebootCtrl,
		NewStatusCtrl,
		NewPutLocationCtrl,
		NewUpdateCtrl,
	),
)
