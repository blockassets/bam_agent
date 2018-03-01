package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/blockassets/bam_agent/controller"
	"github.com/blockassets/bam_agent/fetcher"
	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/service"
	"github.com/blockassets/cgminer_client"
	"github.com/jpillora/overseer"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	// Makefile build
	version        = ""
	interval       time.Duration
	configFileName *string
)

const (
	max24HourInt = 23
	ghUser       = "blockassets"
	ghRepo       = "bam_agent"

	minerHostname = "localhost"
	minerTimeout  = 5 * time.Second
	minerPort     = int64(4028)
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Sometime in the next 24 hours check for update to prevent all machines updating
	// at the same exact time, which could DDOS the network. +1 since rand.Intn is zero based.
	interval = time.Duration(rand.Intn(max24HourInt)+1) * time.Hour

	port := flag.String("port", "1111", "The address to listen on")
	noUpdate := flag.Bool("no-update", false, "Never do any updates. Example: -no-update=true")
	configFileName = flag.String("config", "/etc/bam_agent.json", "configuration file, created if it doesn't exist")
	flag.Parse()

	portStr := fmt.Sprintf(":%s", *port)

	if *noUpdate {
		prog(overseer.State{Address: portStr})
	} else {
		overseerRun(portStr, interval)
	}
}

func overseerRun(port string, interval time.Duration) {
	overseer.Run(overseer.Config{
		Program: prog,
		Address: port,
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

func prog(state overseer.State) {
	log.Printf("%s %s %s %s on port %s", os.Args[0], version, runtime.GOOS, runtime.GOARCH, state.Address)
	if state.Listener != nil {
		log.Printf("Self-update interval: %s", interval)
	}

	cfg, err := LoadAgentConfig(*configFileName)
	if err != nil {
		log.Fatalf("Failed to open configuration: %s\nError: %v\n", *configFileName, err)
		return
	}

	client := minerClient()

	monitorManager := &monitor.Manager{Config: &cfg.Monitor, Client: client}
	monitorManager.StartMonitors()

	startServer(state, client, monitorManager)
}

func startServer(state overseer.State, client *cgminer_client.Client, monitorManager *monitor.Manager) {
	e := echo.New()

	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Must exist here and not as a controller due to issues with rice not supporting nested boxes
	// https://github.com/GeertJohan/go.rice#todo--development  "find boxes in imported packages"
	e.GET("/favicon.ico", echo.WrapHandler(http.FileServer(rice.MustFindBox("static").HTTPBox())))

	controller.Init(e, &controller.Config{Version: version, Client: client, MonitorManager: monitorManager})

	// Start server
	if state.Listener != nil {
		e.Listener = state.Listener
	}
	e.Logger.Fatal(e.Start(state.Address))
}

func minerClient() *cgminer_client.Client {
	port := minerPort

	config, err := service.LoadMinerConfig()
	if err == nil {
		port, err = strconv.ParseInt(config.Path("api-port").Data().(string), 10, 64)
		if err != nil {
			port = minerPort
		}
	}

	return cgminer_client.New(minerHostname, port, minerTimeout)
}
