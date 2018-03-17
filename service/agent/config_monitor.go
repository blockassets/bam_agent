package agent

import (
	"go.uber.org/fx"
)

type ConfigMonitor interface {
	Get() MonitorConfig
	Update(mc MonitorConfig) error
}

type MonitorHelper struct {
	Config
}

func (mh *MonitorHelper) Get() MonitorConfig {
	return mh.Loaded().Monitor
}

func (mh *MonitorHelper) Update(mc MonitorConfig) error {
	return mh.Config.Update("monitor", mc)
}

func NewConfigMonitor(cfg Config) ConfigMonitor {
	return &MonitorHelper{
		Config: cfg,
	}
}

func NewMonitorConfig(cfg ConfigMonitor) MonitorConfig {
	return cfg.Get()
}

var ConfigMonitorModule = fx.Provide(
	NewConfigMonitor,
	NewMonitorConfig,
)
