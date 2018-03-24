package cgminer

import (
	"strconv"

	"go.uber.org/fx"
)

const (
	defaultPort = int64(4028)
)

type ConfigPort interface {
	Get() int64
}

type PortHelper struct {
	Config
}

func (helper *PortHelper) Get() int64 {
	result, ok := helper.Data().Path("api-port").Data().(string)
	if !ok {
		return defaultPort
	}

	val, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return defaultPort
	}
	return val
}

func NewConfigPort(config Config) ConfigPort {
	return &PortHelper{Config: config}
}

var PortModule = fx.Provide(NewConfigPort)
