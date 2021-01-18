package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

func enableMetrics(cctx *cli.Context) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(cctx.String("metrics"), nil)
}
