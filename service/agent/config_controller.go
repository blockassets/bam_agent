package agent

import (
	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

type ConfigController interface {
	Get() ControllerConfig
	Update(cc ControllerConfig) error
}

type ControllerHelper struct {
	Config
}

func (ch *ControllerHelper) Get() ControllerConfig {
	return ch.Loaded().Controller
}

func (ch *ControllerHelper) Update(cc ControllerConfig) error {
	converted, err := jsoniter.Marshal(cc)
	if err != nil {
		return err
	}

	updated, err := gabs.ParseJSON(converted)
	if err != nil {
		return err
	}

	orig := ch.Original()
	orig.Set(updated.Data(), "controller")

	return ch.Save()
}

func NewConfigController(cfg Config) ConfigController {
	return &ControllerHelper{
		Config: cfg,
	}
}

func NewControllerConfig(cfg ConfigController) ControllerConfig {
	return cfg.Get()
}

var ConfigControllerModule = fx.Provide(
	NewConfigController,
	NewControllerConfig,
)