package agent

import (
	"errors"

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
	if mc.Position > 0 && mc.Shelf > 0 {
		return mh.Config.Update("location", mc)
	} else {
		return errors.New("location position and shelf must be greater than 0")
	}
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
