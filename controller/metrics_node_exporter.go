package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/node_exporter/collector"
	"gopkg.in/alecthomas/kingpin.v2"
)

func nodeExporterHandler(filters []string) http.Handler {
	registry := prometheus.NewRegistry()

	nc, err := collector.NewNodeCollector(filters...)
	if err == nil {
		registry.Register(nc)
	} else {
		fmt.Println(err)
	}

	return promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{
			ErrorLog:      nil,
			ErrorHandling: promhttp.ContinueOnError,
		},
	)
}

var collectors = []string{"cpu", "filesystem", "loadavg", "meminfo", "netdev", "netstat", "stat", "time", "uname"}

func NewNodeExporterCtrl() Result {
	return Result{
		Controller: &Controller{
			Path:    "/metrics/node_exporter",
			Methods: []string{http.MethodGet},
			Handler: func() http.Handler {
				// lame workaround due to node_exporter depending on kingpin
				oldArgs := os.Args
				os.Args = []string{""}
				kingpin.Parse()
				os.Args = oldArgs

				return nodeExporterHandler(collectors)
			}(),
		},
	}
}
