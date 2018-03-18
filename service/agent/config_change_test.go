package agent

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/blockassets/bam_agent/tool"
)

// simulate a BAM agent binary update that adds a structure to the default BAM interface
// by saving a previous config file that doesnt have the current monitors in it

const priorConfigVersion = `{
	"cmdLine": {
		"port": ":1111"
	},
	"controller": {
		"reboot": {
			"delay": "5s"
		}
	}
}
`

func TestStructChangeToConfig(t *testing.T) {
	file, err := ioutil.TempFile("/tmp", "agent-config")
	defer os.Remove(file.Name())

	ioutil.WriteFile(file.Name(), []byte(priorConfigVersion), 644)
	file.Close()

	if err != nil {
		t.Fatal(err)
	}

	cfg := NewConfig(tool.CmdLine{
		AgentConfigPath: file.Name(),
	})

	if !cfg.Loaded().Monitor.HighLoad.Enabled {
		t.Fatalf("expected highLoad to be enabled")
	}
}
