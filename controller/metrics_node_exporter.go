package controller

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/node_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

func nodeExporterHandler() http.Handler {
	registry := prometheus.NewRegistry()

	nc, err := collector.NewNodeCollector("cpu", "filesystem", "loadavg", "meminfo", "netdev", "netstat", "stat", "time", "uname")
	if err == nil {
		registry.Register(nc)
	}

	return promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{
			ErrorLog:      nil,
			ErrorHandling: promhttp.ContinueOnError,
		},
	)
}

func NewNodeExporterCtrl() Result {
	return Result{
		Controller: &Controller{
			Path:    "/metrics/node_exporter",
			Methods: []string{http.MethodGet},
			Handler: func() http.Handler {
				// lame
				kingpin.CommandLine.Terminate(nil)
				kingpin.Parse()

				return nodeExporterHandler()
			}(),
		},
	}
}
