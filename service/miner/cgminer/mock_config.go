package cgminer

import (
	"log"

	"github.com/Jeffail/gabs"
)

const (
	MockDefaultConfigFile = `{
 "HttpPort": "80",
 "pool1": "",
 "pool2": "",
 "pool3": "",
 "failover-only": true,
 "no-submit-stale": true,
 "api-listen": true,
 "api-port": "4028",
 "api-allow": "W:0/0",
 "volt": "2",
 "username": "admin",
 "password": "bw.com",
 "language": "ch",
 "chipNumber": "36",
 "frequency": "684",
 "autoFrequency": true,
 "autoGetJobTimeOut": true,
 "frequencySet": "384_30|450_30|480_30|540_30|576_30|600_30|612_30|625_30|636_30|648_30|660_29|672_29|684_28|700_28|720_28|744_28|756_28|768_28|800_28|912_28|1020_28",
 "fanSet": "30_1000|34_2000|38_3000|42_4000|46_5000|50_6000",
 "autoNet": false,
 "ip": "10.10.0.11",
 "mask": "255.255.252.0",
 "gateway": "10.10.0.1",
 "dns": "8.8.8.8",
 "debug": true,
 "packet": true,
 "botelv": true,
 "board_reset_waittime": "14",
 "mcu_reset_waittime": "0",
 "invalid_cnt": "30",
 "scanwork_sleeptime": "4",
 "board_reenable_waittime": "60",
 "temp_threshold": "80",
 "task_interval": "350",
 "start_voltage": "6000",
 "running_voltage1": "5650",
 "running_voltage2": "5650",
 "running_voltage3": "5650",
 "fengru": "5000",
 "fengchu": "5000"
}`
)

// Type insurance
var _ Config = &MockConfig{}

type MockConfig struct {
	data      string
	container *gabs.Container
	Saved     string
}

func (mc *MockConfig) Load() error {
	container, err := gabs.ParseJSON([]byte(mc.data))
	if err != nil {
		return err
	}
	mc.container = container
	return nil
}

func (mc *MockConfig) Save() error {
	mc.Saved = string(mc.container.BytesIndent("", "\t"))
	return nil
}

func (mc *MockConfig) Data() *gabs.Container {
	return mc.container
}

func NewMockConfig(data string) MockConfig {
	if len(data) == 0 {
		data = MockDefaultConfigFile
	}
	mc := MockConfig{
		data: data,
	}
	err := mc.Load()
	if err != nil {
		log.Fatalf("could not load MockConfig %s", data)
	}
	return mc
}
