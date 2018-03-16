package miner

import (
	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

type PoolAddresses struct {
	Pool1 string `json:"pool1"`
	Pool2 string `json:"pool2"`
	Pool3 string `json:"pool3"`
}

type ConfigPools interface {
	Parse(data []byte) (*PoolAddresses, error)
	Get() (*PoolAddresses, error)
	Save(pools *PoolAddresses) error
}

type PoolHelper struct {
	Config
}

func NewPoolHelper(config Config) ConfigPools {
	return &PoolHelper{config}
}

var PoolModule = fx.Provide(NewPoolHelper)

func (helper PoolHelper) Parse(data []byte) (*PoolAddresses, error) {
	pools := &PoolAddresses{}
	err := jsoniter.Unmarshal(data, pools)
	if err != nil {
		return nil, err
	}
	return pools, nil
}

func (helper PoolHelper) Save(pools *PoolAddresses) error {
	c := helper.Data()
	c.Set(pools.Pool1, "pool1")
	c.Set(pools.Pool2, "pool2")
	c.Set(pools.Pool3, "pool3")

	return helper.Config.Save()
}

func (helper PoolHelper) Get() (*PoolAddresses, error) {
	c := helper.Data()
	pool1, _ := c.Path("pool1").Data().(string)
	pool2, _ := c.Path("pool2").Data().(string)
	pool3, _ := c.Path("pool3").Data().(string)

	return &PoolAddresses{
		Pool1: pool1,
		Pool2: pool2,
		Pool3: pool3,
	}, nil
}
