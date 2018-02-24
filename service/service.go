package service

type Commands interface {
	CgmQuit() error
	Reboot()
}

type Command struct {
}
