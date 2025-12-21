package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func startMetricsEndpoint() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Prometheus metrics available at /metrics")
	log.Fatal(http.ListenAndServe(":9091", nil))
}
