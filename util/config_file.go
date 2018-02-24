package util

import (
	"io/ioutil"
	"os"
)

var defaultConfig = []byte(`{ 
	"monitor": {
		"miner_quit":{	
			"enabled":true,
			"period_secs":15,
			"initial_range_secs":5
		}, 
		"system_reboot":{	
			"enabled":true,
			"period_secs":15,
			"initial_range_secs":5
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
