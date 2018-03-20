package agent

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

const priorConfigVersion = `{
	"cmdLine": {
		"port": ":2222"
	},
	"controller": {
		"reboot": {
			"delay": "10s"
		}
	}
}
`

const failingConfig = `
{
  "cmdLine": {
    "port": ":1111"
  },
  "controller": {
    "reboot": {
      "delay": "5s"
    }
  },
  "monitor": {
    "load": {
      "period_secs": 60,
      "high_load_Mark": 5,
      "enabled": true
    }
  }
}
`

func defaultConfig() []byte {
	data, _ := ioutil.ReadFile("../../conf/bam_agent.json")
	return data
}

func TestNewConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "agent-config")
	defer file.Close()
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	data := defaultConfig()
	if len(data) == 0 {
		t.Fatal("no data!")
	}

	// Test loading an empty file
	cfg := NewConfig(tool.CmdLine{
		AgentConfigPath: file.Name(),
	},
		data,
	)

	if !cfg.Loaded().Monitor.HighLoad.Enabled {
		t.Fatalf("expected highLoad to be enabled")
	}

	fileData, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(fileData) == 0 {
		t.Fatalf("expected file to have data, got %v", fileData)
	}

	fileConfig := &FileConfig{}
	err = jsoniter.Unmarshal(fileData, fileConfig)
	if err != nil {
		t.Fatal(err)
	}

	if fileConfig.Controller.Reboot.Delay != time.Duration(5)*time.Second {
		t.Fatalf("expected 5s and got %v", fileConfig.Controller.Reboot.Delay)
	}
}

func TestLoadingOldConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "agent-config")
	defer os.Remove(file.Name())

	// Write out an old version of a file
	ioutil.WriteFile(file.Name(), []byte(failingConfig), 644)
	file.Close()

	if err != nil {
		t.Fatal(err)
	}

	cfg := NewConfig(tool.CmdLine{
		AgentConfigPath: file.Name(),
	}, defaultConfig())

	if !cfg.Loaded().Monitor.HighLoad.Enabled {
		t.Fatalf("expected highLoad to be enabled")
	}

	// We should have saved the file as part of the load
	fileData, err := ioutil.ReadFile(file.Name())
	if !strings.Contains(string(fileData), "highLoad") {
		t.Fatalf("expected file data to have updated port number")
	}
}

// simulate a BAM agent binary update that adds a structure to the default BAM interface
// by saving a previous config file that doesnt have the current monitors in it
func TestStructChangeToConfig(t *testing.T) {
	file, err := ioutil.TempFile("", "agent-config")
	defer os.Remove(file.Name())

	// Write out an old version of a file
	ioutil.WriteFile(file.Name(), []byte(priorConfigVersion), 644)
	file.Close()

	if err != nil {
		t.Fatal(err)
	}

	cfg := NewConfig(tool.CmdLine{
		AgentConfigPath: file.Name(),
	}, defaultConfig())

	if !cfg.Loaded().Monitor.HighLoad.Enabled {
		t.Fatalf("expected highLoad to be enabled")
	}

	// Test that the file is merged on top of the defaults
	if cfg.Loaded().Controller.Reboot.Delay != time.Duration(10)*time.Second {
		t.Fatalf("expected 10s controller reboot delay, got %v", cfg.Loaded().Controller.Reboot.Delay)
	}

	if cfg.Loaded().CmdLine.Port != ":2222" {
		t.Fatalf("expected cmdline port to be :2222, got: %s", cfg.Loaded().CmdLine.Port)
	}

	// We should have saved the file as part of the load
	fileData, err := ioutil.ReadFile(file.Name())
	if !strings.Contains(string(fileData), ":2222") {
		t.Fatalf("expected file data to have updated port number")
	}
}
