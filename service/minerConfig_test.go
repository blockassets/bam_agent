package service

import (
	"encoding/json"
	"testing"

	"github.com/Jeffail/gabs"
)

var updatePoolsTests = []struct {
	configFile  string
	poolsIn     string
	shouldErr   bool
	expectedOut string
}{
	// happy path
	{
		configFile:`{
	"api-allow": "W:0/0",
	"api-listen": true,
	"api-port": "4028",
	"autoFrequency": true,
	"autoGetJobTimeOut": true,
	"autoNet": false,
	"board_reenable_waittime": "60",
	"board_reset_waittime": "14",
	"botelv": true,
	"chipNumber": "36",
	"debug": true,
	"dns": "8.8.8.8",
	"failover-only": true,
	"fanSet": "30_1000|34_2000|38_3000|42_4000|46_5000|50_6000",
	"fengchu": "5000",
	"fengru": "5000",
	"frequency": "684",
	"frequencySet": "384_30|450_30|480_30|540_30|576_30|600_30|612_30|625_30|636_30|648_30|660_29|672_29|684_28|700_28|720_28|744_28|756_28|768_28|800_28|912_28|1020_28",
	"gateway": "",
	"invalid_cnt": "30",
	"ip": "",
	"language": "ch",
	"mask": "",
	"mcu_reset_waittime": "0",
	"no-submit-stale": true,
	"packet": true,
	"password": "bw.com",
	"pool1": "111.2.3.4",
	"pool2": "111.3.4.5",
	"pool3": "111.4.5.6",
	"running_voltage1": "5650",
	"running_voltage2": "5650",
	"running_voltage3": "5650",
	"scanwork_sleeptime": "4",
	"start_voltage": "6000",
	"task_interval": "350",
	"temp_threshold": "80",
	"username": "admin",
	"volt": "2"
}`,
	poolsIn: `{"pool1": "333.2.3.4", "pool2": "333.3.4.5", "pool3": "333.4.5.6"}`,
	expectedOut: `{
	"api-allow": "W:0/0",
	"api-listen": true,
	"api-port": "4028",
	"autoFrequency": true,
	"autoGetJobTimeOut": true,
	"autoNet": false,
	"board_reenable_waittime": "60",
	"board_reset_waittime": "14",
	"botelv": true,
	"chipNumber": "36",
	"debug": true,
	"dns": "8.8.8.8",
	"failover-only": true,
	"fanSet": "30_1000|34_2000|38_3000|42_4000|46_5000|50_6000",
	"fengchu": "5000",
	"fengru": "5000",
	"frequency": "684",
	"frequencySet": "384_30|450_30|480_30|540_30|576_30|600_30|612_30|625_30|636_30|648_30|660_29|672_29|684_28|700_28|720_28|744_28|756_28|768_28|800_28|912_28|1020_28",
	"gateway": "",
	"invalid_cnt": "30",
	"ip": "",
	"language": "ch",
	"mask": "",
	"mcu_reset_waittime": "0",
	"no-submit-stale": true,
	"packet": true,
	"password": "bw.com",
	"pool1": "333.2.3.4",
	"pool2": "333.3.4.5",
	"pool3": "333.4.5.6",
	"running_voltage1": "5650",
	"running_voltage2": "5650",
	"running_voltage3": "5650",
	"scanwork_sleeptime": "4",
	"start_voltage": "6000",
	"task_interval": "350",
	"temp_threshold": "80",
	"username": "admin",
	"volt": "2"
}`},
}

func TestMutateConfig(t *testing.T) {
	for index, tt := range updatePoolsTests {
		jsonConfig, err := gabs.ParseJSON([]byte(tt.configFile))
		if err != nil {
			t.Error(err)
		}

		pools := &PoolAddresses{}
		err = json.Unmarshal([]byte(tt.poolsIn), pools)
		if err != nil {
			t.Error(err)
		}

		mutated := mutateConfig(pools, jsonConfig)
		buf := string(mutated)
		if tt.expectedOut != buf {
			t.Errorf("Test Index: %v: Expected:\n%s\nGot:\n%s\n ", index, tt.expectedOut, buf)
		}
	}
}

func TestUpdatePools(t *testing.T) {
	err := UpdatePools(nil)
	if err == nil {
		t.Error("Should have had an error on nil input")
	}

	err = UpdatePools([]byte("{ this is bad json }"))
	if err == nil {
		t.Error("Should have had an error on bad input")
	}
}
