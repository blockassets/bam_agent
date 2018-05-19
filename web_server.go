package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/blockassets/bam_agent/controller"
	"github.com/blockassets/bam_agent/tool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lookfirst/overseer"
	"go.uber.org/fx"
)

type WebServer struct {
	echo          *echo.Echo
	ctrl          controller.Manager
	state         overseer.State
	staticRiceBox tool.StaticRiceBox
}

func NewWebServer(e *echo.Echo, ctrl controller.Manager, state overseer.State, staticRiceBox tool.StaticRiceBox) *WebServer {
	return &WebServer{
		echo:          e,
		ctrl:          ctrl,
		state:         state,
		staticRiceBox: staticRiceBox,
	}
}

func (server *WebServer) Start() {
	go func() {
		go run(server.echo, server.state, server.staticRiceBox)

		// Blocks until we receive a shutdown notice
		<-server.state.GracefulShutdown

		stop(server)
	}()
}

func (server *WebServer) Stop() {
	server.state.GracefulShutdown <- true
}

func stop(server *WebServer) {
	// After 10 seconds we gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.echo.Shutdown(ctx); err != nil {
		server.echo.Logger.Fatal(err)
	} else {
		log.Println("Shutdown")
	}
}

func run(e *echo.Echo, state overseer.State, staticRiceBox tool.StaticRiceBox) {
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/favicon.ico", echo.WrapHandler(http.FileServer((*staticRiceBox).HTTPBox())))

	// Start server
	if state.Listener != nil {
		e.Listener = state.Listener
	}

	e.Logger.Fatal(e.Start(state.Address))
}

func NewEcho() *echo.Echo {
	return echo.New()
}

var WebServerModule = fx.Options(
	controller.Module,

	fx.Provide(
		NewEcho,
		NewWebServer,
	),
)
