package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/blockassets/bam_agent/controller"
	"github.com/jpillora/overseer"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/fx"
)

type WebServer struct {
	echo  *echo.Echo
	ctrl  controller.Manager
	state overseer.State
}

func NewWebServer(e *echo.Echo, ctrl controller.Manager, state overseer.State) *WebServer {
	return &WebServer{
		echo:  e,
		ctrl:  ctrl,
		state: state,
	}
}

func (server *WebServer) Start() {
	go func() {
		go run(server.echo, server.state)

		// Blocks until we receive a shutdown notice
		<-server.state.GracefulShutdown

		// After 10 seconds we gracefully shutdown the server
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.echo.Shutdown(ctx); err != nil {
			server.echo.Logger.Fatal(err)
		} else {
			log.Println("Shutdown")
		}
	}()
}

func run(e *echo.Echo, state overseer.State) {
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Must exist here and not as a controller due to issues with rice not supporting nested boxes
	// https://github.com/GeertJohan/go.rice#todo--development  "find boxes in imported packages"
	e.GET("/favicon.ico", echo.WrapHandler(http.FileServer(rice.MustFindBox("static").HTTPBox())))

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
