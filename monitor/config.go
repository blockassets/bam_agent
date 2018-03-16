package monitor

import (
	"github.com/blockassets/bam_agent/service/agent"
	"go.uber.org/fx"
)

// Copy the config from the agent data source in order to keep the dependency between
// the config file (Entity) and the monitor layer clean via a DTO. Managed through the
// dependency injection so it is all cleanly done.
var ConfigProviders = fx.Options(
	fx.Provide(func(cfg agent.Config) AcceptedConfig {
		value := cfg.Monitor.AcceptedShares
		return AcceptedConfig{
			Enabled: value.Enabled,
			Period:  value.Period,
		}
	}),

	fx.Provide(func(cfg agent.Config) HighLoadConfig {
		value := cfg.Monitor.HighLoad
		return HighLoadConfig{
			Enabled:      value.Enabled,
			Period:       value.Period,
			HighLoadMark: value.HighLoadMark,
		}
	}),

	fx.Provide(func(cfg agent.Config) HighTempConfig {
		value := cfg.Monitor.HighTemp
		return HighTempConfig{
			Enabled:  value.Enabled,
			Period:   value.Period,
			HighTemp: value.HighTemp,
		}
	}),

	fx.Provide(func(cfg agent.Config) CGMQuitConfig {
		value := cfg.Monitor.CGMQuit
		return CGMQuitConfig{
			Enabled: value.Enabled,
			Period:  value.Period.Duration,
		}
	}),

	fx.Provide(func(cfg agent.Config) RebootConfig {
		value := cfg.Monitor.Reboot
		return RebootConfig{
			Enabled: value.Enabled,
			Period:  value.Period.Duration,
		}
	}),
)
