package miner

import (
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

type NetworkData struct {
	IPAddress string `json:"ip"`
	Netmask   string `json:"mask"`
	Gateway   string `json:"gateway"`
	Dns       string `json:"dns"`
}

type ConfigNetwork interface {
	Parse(data []byte) (*NetworkData, error)
	Get() (*NetworkData, error)
	Save(data *NetworkData) error
}

type NetworkHelper struct {
	Config
}

func (helper NetworkHelper) Parse(data []byte) (*NetworkData, error) {
	net := &NetworkData{}
	err := jsoniter.Unmarshal(data, net)
	if err != nil {
		return nil, err
	}
	return net, nil
}

func (helper NetworkHelper) Get() (*NetworkData, error) {
	c := helper.Data()
	addr, _ := c.Path("ip").Data().(string)
	netmask, _ := c.Path("mask").Data().(string)
	gateway, _ := c.Path("gateway").Data().(string)
	dns, _ := c.Path("dns").Data().(string)

	return &NetworkData{
		IPAddress: addr,
		Netmask:   netmask,
		Gateway:   gateway,
		Dns:       dns,
	}, nil
}

func (helper NetworkHelper) Save(data *NetworkData) error {
	c := helper.Data()
	if len(data.IPAddress) == 0 {
		c.Set(true, "autoNet")
	}
	c.Set(data.IPAddress, "address")
	c.Set(data.Netmask, "netmask")
	c.Set(data.Gateway, "gateway")
	c.Set(data.Dns, "dns")

	return helper.Config.Save()
}

func NewConfigNetwork(config Config) ConfigNetwork {
	return &NetworkHelper{config}
}

var NetworkModule = fx.Provide(NewConfigNetwork)
