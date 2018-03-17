package miner

import (
	"strconv"
	"sync"

	"github.com/json-iterator/go"
	"go.uber.org/fx"
)

const (
	defaultFrequency = int64(684)
)

type FrequencyData struct {
	Frequency int64 `json:"frequency,string"`
}

type ConfigFrequency interface {
	Get() int64
	Save(int64) error
	Parse(data []byte) (*FrequencyData, error)
}

type FrequencyHelper struct {
	Config
	sync.Mutex
}

func (helper FrequencyHelper) Get() int64 {
	result, ok := helper.Data().Path("frequency").Data().(string)
	if !ok {
		return defaultFrequency
	}

	val, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return defaultFrequency
	}
	return val
}

func (helper FrequencyHelper) Save(frequency int64) error {
	helper.Lock()
	defer helper.Unlock()
	c := helper.Config.Data()
	c.Set(strconv.FormatInt(frequency, 10), "frequency")
	return helper.Config.Save()
}

func (helper FrequencyHelper) Parse(data []byte) (*FrequencyData, error) {
	mf := &FrequencyData{}
	err := jsoniter.Unmarshal(data, mf)
	if err != nil {
		return nil, err
	}

	return mf, nil
}

func NewConfigFrequency(config Config) ConfigFrequency {
	return &FrequencyHelper{Config: config}
}

var FrequencyModule = fx.Provide(NewConfigFrequency)
