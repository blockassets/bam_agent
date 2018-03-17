package agent

import (
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
	return ch.Config.Update("controller", cc)
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
