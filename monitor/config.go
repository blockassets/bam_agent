package monitor

import (
	"github.com/blockassets/bam_agent/service/agent"
	"go.uber.org/fx"
)

// Copy the config from the agent data source in order to keep the dependency between
// the config file (Entity) and the monitor layer clean via a DTO. Managed through the
// dependency injection so it is all cleanly done.
var ConfigProviders = fx.Options(
	fx.Provide(func(cfg agent.MonitorConfig) AcceptedConfig {
		value := cfg.AcceptedShares
		return AcceptedConfig{
			Enabled: value.Enabled,
			Period:  value.Period,
		}
	}),

	fx.Provide(func(cfg agent.MonitorConfig) HighLoadConfig {
		value := cfg.HighLoad
		return HighLoadConfig{
			Enabled:      value.Enabled,
			Period:       value.Period,
			HighLoadMark: value.HighLoadMark,
		}
	}),

	fx.Provide(func(cfg agent.MonitorConfig) HighTempConfig {
		value := cfg.HighTemp
		return HighTempConfig{
			Enabled:  value.Enabled,
			Period:   value.Period,
			HighTemp: value.HighTemp,
		}
	}),

	fx.Provide(func(cfg agent.MonitorConfig) CGMQuitConfig {
		value := cfg.CGMQuit
		return CGMQuitConfig{
			Enabled: value.Enabled,
			Period:  value.Period.Duration,
		}
	}),

	fx.Provide(func(cfg agent.MonitorConfig) RebootConfig {
		value := cfg.Reboot
		return RebootConfig{
			Enabled: value.Enabled,
			Period:  value.Period.Duration,
		}
	}),

	fx.Provide(func(cfg agent.MonitorConfig) LowMemoryConfig {
		value := cfg.LowMemory
		return LowMemoryConfig{
			Enabled:   value.Enabled,
			Period:    value.Period,
			LowMemory: value.LowMemory,
		}
	}),

	fx.Provide(func(cfg agent.MonitorConfig) NtpdateConfig {
		value := cfg.Ntpdate
		return NtpdateConfig{
			Enabled: value.Enabled,
			Period:  value.Period,
		}
	}),
)
