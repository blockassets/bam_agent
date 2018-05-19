package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/blockassets/bam_agent/fetcher"
	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service/agent"
	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/bam_agent/service/miner/cgminer"
	"github.com/blockassets/bam_agent/service/os"
	"github.com/blockassets/bam_agent/tool"
	"github.com/lookfirst/overseer"
	"go.uber.org/fx"
)

var (
	cmdLine        tool.CmdLine
	monitorManager monitor.Manager
	webManager     *WebServer
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
	webManager = ws
	ws.Start()
}

func monitors(mgr monitor.Manager) {
	monitorManager = mgr
	mgr.Start()
}

func program(state overseer.State) {
	cmdLineProvider := fx.Provide(func() tool.CmdLine {
		return cmdLine
	})

	stateProvider := fx.Provide(func() overseer.State {
		return state
	})

	staticRiceBox := fx.Provide(func() tool.StaticRiceBox {
		return rice.MustFindBox("static")
	})

	confRiceBox := fx.Provide(func() tool.ConfRiceBox {
		return rice.MustFindBox("conf")
	})

	app := fx.New(
		cmdLineProvider,
		stateProvider,

		staticRiceBox,
		confRiceBox,

		agent.ConfigModule,
		agent.VersionModule,

		cgminer.ConfigModule,
		miner.ClientModule,
		miner.VersionModule,

		os.MemInfoModule,
		os.MinerModule,
		os.NetInfoModule,
		os.NetworkingModule,
		os.NtpdateModule,
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
		program(overseer.State{
			Address:          cmdLine.Port,
			GracefulShutdown: make(chan bool, 1),
		})
	} else {
		overseerRun(cmdLine.Port)
	}
}

func overseerRun(port string) {
	interval := time.Duration(rand.Intn(5)+1) * time.Hour // within the next 6 hours

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
		// Try to prevent /reboot from being called during an upgrade
		PreUpgrade: func(tempBinaryPath string) error {
			if monitorManager != nil {
				monitorManager.Stop()
			}

			if webManager != nil {
				webManager.Stop()
			}
			return nil
		},
	})
}
