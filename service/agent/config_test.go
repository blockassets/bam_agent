package agent

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/blockassets/bam_agent/tool"
	"github.com/json-iterator/go"
)

func TestNewConfig(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "agent-config")
	defer file.Close()
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	cfg := NewConfig(tool.CmdLine{
		AgentConfigPath: file.Name(),
	})

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

	if fileConfig.Controller.Reboot.Delay != time.Duration(5) * time.Second {
		t.Fatalf("expected 5s and got %v", fileConfig.Controller.Reboot.Delay)
	}
}
