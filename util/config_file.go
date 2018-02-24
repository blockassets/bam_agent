package util

import (
	"io/ioutil"
	"os"
)
// day == 86400 seconds, hour == 3600 seconds
// 3 day == 259200
var defaultConfig = []byte(`{ 
	"monitor": {
		"miner_quit":{	
			"enabled":false,
			"period_secs":86400,
			"initial_range_secs":3600
		}, 
		"system_reboot":{	
			"enabled":false,
			"period_secs":259200,
			"initial_range_secs":3600
		},
		"system_load":{
			"enabled":true,
			"period_secs":60,
			"initial_range_secs":1
		}
	}
}`)

type ConfigFile struct {
	configFileName string
	configBuffer   []byte
}

func InitialiseConfigFile(configFile string) (*ConfigFile, error) {
	c := &ConfigFile{}
	c.configFileName = configFile

	f, err := os.OpenFile(c.configFileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err == nil {
		// the file didnt exist before, so write out the defualt
		f.Write(defaultConfig)
		f.Close()
		c.configBuffer = defaultConfig
	} else {
		if os.IsExist(err) {
			c.configBuffer, err = ioutil.ReadFile(c.configFileName)
			if err != nil {
				return c, err
			}
		}else {
			return c,err
		}
	}
	return c, nil
}

func (c *ConfigFile) GetConfigBuf() []byte {
	return c.configBuffer
}
