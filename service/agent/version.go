package agent

import "go.uber.org/fx"

var (
	// Makefile build
	version = ""
)

type Version struct {
	V string
}

func NewVersion() Version {
	if len(version) > 0 {
		return Version{V: version}
	}
	return Version{}
}

var VersionModule = fx.Provide(NewVersion)
