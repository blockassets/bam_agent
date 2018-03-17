package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/blockassets/bam_agent/fetcher"
	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/agent"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/jpillora/overseer"
	"go.uber.org/fx"
)

var (
	cmdLine tool.CmdLine
)

const (
	ghUser = "blockassets"
	ghRepo = "bam_agent"
)

func setup(version agent.Version) {
	rand.Seed(time.Now().UTC().UnixNano())
	log.Printf("Agent version: %s ", version.V)
}

func webServer(ws *WebServer) {
	ws.Start()
}

func monitors(mgr monitor.Manager) {
	mgr.Start()
}

func program(state overseer.State) {
	cmdLineProvider := fx.Provide(func() tool.CmdLine {
		return cmdLine
	})

	stateProvider := fx.Provide(func() overseer.State {
		return state
	})

	app := fx.New(
		cmdLineProvider,
		stateProvider,

		agent.ConfigModule,
		agent.VersionModule,

		miner.ConfigModule,
		miner.ClientModule,
		miner.VersionModule,

		os.MinerModule,
		os.NetInfoModule,
		os.NetworkingModule,
		os.RebootModule,
		os.StatRetrieverModule,
		os.UptimeModule,

		monitor.Module,
		WebServerModule,

		fx.Invoke(setup, webServer, monitors),
	)

	app.Run()
}

/*
	main() gets called 2x when we use overseer. It first starts up a master process and then
	starts the app again as a child process. This means that the command line arguments need to
	be parsed twice. So, we cache them and then make them available to fx injection in the program() function
	by creating a provider for them.
*/
func main() {
	cmdLine = tool.NewCmdLine()

	if cmdLine.NoUpdate {
		program(overseer.State{Address: cmdLine.Port})
	} else {
		overseerRun(cmdLine.Port)
	}
}

func overseerRun(port string) {
	interval := time.Duration(rand.Intn(23)+1) * time.Hour // within the next 24 hours

	overseer.Run(overseer.Config{
		Debug:     true,
		NoRestart: true, // We allow the OS to restart things
		Program:   program,
		Address:   port,
		// The default is to check on startup, but we really just want to check in the next interval
		// in order to prevent DDOS'ing the whole network if we restart all machines. I copied the overseer
		// version of the github fetcher into this project and modified the logic there.
		Fetcher: &fetcher.Github{
			User:     ghUser,
			Repo:     ghRepo,
			Interval: interval,
		},
	})
}
