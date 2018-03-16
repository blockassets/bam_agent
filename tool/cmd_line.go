package tool

import (
	"flag"
	"fmt"
)

/*
	Contains all the variables passed in on the command line.
*/
type CmdLine struct {
	Port            string `json:"port"`
	AgentConfigPath string `json:"config-agent"`
	MinerConfigPath string `json:"config-miner"`
}

func NewCmdLine() CmdLine {
	//noUpdate := flag.Bool("no-update", false, "Never do any updates. Example: -no-update=true")
	port := flag.String("port", "1111", "The address to listen on")
	agentConfigPath := flag.String("config-agent", "/etc/bam_agent.json", "Agent configuration file, created if it doesn't exist")
	minerConfigPath := flag.String("config-miner", "/usr/app/conf.default", "Miner configuration file")
	flag.Parse()

	portStr := fmt.Sprintf(":%s", *port)

	cmdLine := CmdLine{
		Port:            portStr,
		AgentConfigPath: *agentConfigPath,
		MinerConfigPath: *minerConfigPath,
	}

	return cmdLine
}
