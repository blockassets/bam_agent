package service

import (
	"testing"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

const (
	inputConfig = `{
	"api-allow": "W:0/0",
	"api-listen": true,
	"api-port": "4028",
	"autoFrequency": true,
	"autoGetJobTimeOut": true,
	"autoNet": true,
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
}`
	poolsToMutate      = `{"pool1": "333.2.3.4", "pool2": "333.3.4.5", "pool3": "333.4.5.6"}`
	expectedPoolMutate = `{
	"api-allow": "W:0/0",
	"api-listen": true,
	"api-port": "4028",
	"autoFrequency": true,
	"autoGetJobTimeOut": true,
	"autoNet": true,
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
}`
	staticIPAddressesToMutate = `{ "address": "1.2.3.4", "netmask": "5.6.7.8", "gateway": "9.10.11.12", "dns": "13.13.13.13"}`
	expectedStaticIPMutate    = `{
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
	"dns": "13.13.13.13",
	"failover-only": true,
	"fanSet": "30_1000|34_2000|38_3000|42_4000|46_5000|50_6000",
	"fengchu": "5000",
	"fengru": "5000",
	"frequency": "684",
	"frequencySet": "384_30|450_30|480_30|540_30|576_30|600_30|612_30|625_30|636_30|648_30|660_29|672_29|684_28|700_28|720_28|744_28|756_28|768_28|800_28|912_28|1020_28",
	"gateway": "9.10.11.12",
	"invalid_cnt": "30",
	"ip": "1.2.3.4",
	"language": "ch",
	"mask": "5.6.7.8",
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
}`
	expectedDHCPMutate = `{
	"api-allow": "W:0/0",
	"api-listen": true,
	"api-port": "4028",
	"autoFrequency": true,
	"autoGetJobTimeOut": true,
	"autoNet": true,
	"board_reenable_waittime": "60",
	"board_reset_waittime": "14",
	"botelv": true,
	"chipNumber": "36",
	"debug": true,
	"dns": "",
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
}`
)

func TestMutatePools(t *testing.T) {
	jsonConfig, err := gabs.ParseJSON([]byte(inputConfig))
	if err != nil {
		t.Error(err)
	}

	pools := &PoolAddresses{}
	err = jsoniter.Unmarshal([]byte(poolsToMutate), pools)
	if err != nil {
		t.Error(err)
	}

	mutated := mutatePools(pools, jsonConfig)
	buf := string(mutated)
	if expectedPoolMutate != buf {
		t.Errorf("Expected:\n%s\nGot:\n%s\n ", expectedPoolMutate, buf)
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

func TestConfigFilePath(t *testing.T) {
	if configFilePath != "/usr/app/conf.default" {
		t.Errorf("Expected config file path to be /usr/app/conf.default and got %s", configFilePath)
	}
}

func TestUpdateStaticNetConfig(t *testing.T) {
	err := UpdateStaticNetConfig(nil)
	if err == nil {
		t.Error("Should have had an error on nil input")
	}

	err = UpdateStaticNetConfig([]byte("{ this is bad json }"))
	if err == nil {
		t.Error("Should have had an error on bad input")
	}
}

func TestMutateStaticNetConfig(t *testing.T) {
	jsonConfig, err := gabs.ParseJSON([]byte(inputConfig))
	if err != nil {
		t.Error(err)
	}

	netConfig := &StaticNetConfig{}
	err = jsoniter.Unmarshal([]byte(staticIPAddressesToMutate), netConfig)
	if err != nil {
		t.Error(err)
	}

	mutated := mutateStaticNetConfig(netConfig, jsonConfig)
	buf := string(mutated)
	if expectedStaticIPMutate != buf {
		t.Errorf("Expected:\n%s\nGot:\n%s\n ", expectedStaticIPMutate, buf)
	}
}

func TestMutateDHCPNetConfig(t *testing.T) {
	// Use the expected Static configuration as input...
	jsonConfig, err := gabs.ParseJSON([]byte(expectedStaticIPMutate))
	if err != nil {
		t.Error(err)
	}

	mutated := mutateDHCPNetConfig(jsonConfig)
	buf := string(mutated)
	if expectedDHCPMutate != buf {
		t.Errorf("Expected:\n%s\nGot:\n%s\n ", expectedDHCPMutate, buf)
	}
}
