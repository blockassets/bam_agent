package agent

import (
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

type ConfigLocation interface {
	Parse(data []byte) (*LocationConfig, error)
	Get() LocationConfig
	Update(mc LocationConfig) error
}

type LocationHelper struct {
	Config
}

func (mh *LocationHelper) Parse(data []byte) (*LocationConfig, error) {
	locationCfg := &LocationConfig{}
	err := jsoniter.Unmarshal(data, locationCfg)
	if err != nil {
		return nil, err
	}
	return locationCfg, nil
}

func (mh *LocationHelper) Get() LocationConfig {
	return mh.Loaded().Location
}

func (mh *LocationHelper) Update(mc LocationConfig) error {
	return mh.Config.Update("location", mc)
}

func NewConfigLocation(cfg Config) ConfigLocation {
	return &LocationHelper{
		Config: cfg,
	}
}

func NewLocationConfig(cfg ConfigLocation) LocationConfig {
	return cfg.Get()
}

var ConfigLocationModule = fx.Provide(
	NewConfigLocation,
	NewLocationConfig,
)
