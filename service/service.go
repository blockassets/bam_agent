package service

import (
	"io"
)

type Commands interface {
	CgmQuit() error
	Reboot()
	UpdatePools(poolsAsJson io.ReadCloser, configFilePath string) error
}

type Command struct {
}
