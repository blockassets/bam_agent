package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/blockassets/bam_agent/util"
	"github.com/blockassets/bam_agent/monitor"
	"github.com/blockassets/bam_agent/controller"
)

var (
	// Makefile build
	version = ""
)

const (
	theBAMConfigFile = "./etc/bam_agent.json"
)

func main() {
	log.Printf("%s %s", os.Args[0], version)
	cfg, err := util.InitialiseConfigFile(theBAMConfigFile)
	if err != nil {
		log.Printf("Can't initialize Configuration: %v", err)
		return
	}
	monitor.StartMonitors(cfg)
	startServer()
}

func startServer() {
	port := flag.String("port", "1111", "The address to listen on.")
	flag.Parse()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/favicon.ico", echo.WrapHandler(http.FileServer(rice.MustFindBox("static").HTTPBox())))

	controller.Init(e)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", *port)))
}
