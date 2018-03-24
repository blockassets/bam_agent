package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/service/miner"
	"github.com/blockassets/cgminer_client"
	"github.com/blockassets/cgminer_exporter/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewCgMinerExporterCtrl(client cgminer_client.Client, version miner.Version) Result {
	return Result{
		Controller: &Controller{
			Path:    "/metrics/cgminer_exporter",
			Methods: []string{http.MethodGet},
			Handler: func() http.Handler {
				exporter := exporter.NewExporter(client, version.V)

				registry := prometheus.NewRegistry()
				registry.Register(exporter)

				return promhttp.HandlerFor(registry,
					promhttp.HandlerOpts{
						ErrorLog:      nil,
						ErrorHandling: promhttp.ContinueOnError,
					},
				)
			}(),
		},
	}
}
