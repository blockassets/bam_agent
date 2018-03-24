package controller

import (
	"net/http"

	"github.com/blockassets/node_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

func makeHandler() http.Handler {
	registry := prometheus.NewRegistry()
	nc, err := collector.NewNodeCollector()
	if err != nil {
		log.Fatalf("Couldn't create collector: %s", err)
	}

	registry.MustRegister(nc)

	return promhttp.HandlerFor(registry,
		promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		},
	)
}

func NewNodeExporterCtrl() Result {
	return Result{
		Controller: &Controller{
			Path:    "/metrics/node_exporter",
			Methods: []string{http.MethodGet},
			Handler: makeHandler(),
		},
	}
}
