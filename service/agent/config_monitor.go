package agent

import (
	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
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
	converted, err := jsoniter.Marshal(mc)
	if err != nil {
		return err
	}

	updated, err := gabs.ParseJSON(converted)
	if err != nil {
		return err
	}

	orig := mh.Original()
	orig.Set(updated.Data(), "monitor")

	return mh.Save()
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