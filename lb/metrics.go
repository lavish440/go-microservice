package main

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Load Balancer Metrics
var (
	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadbalancer_active_connections",
			Help: "Current number of active TCP connections",
		},
	)

	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadbalancer_requests_total",
			Help: "Total number of requests served by the load balancer",
		},
		[]string{"backend"},
	)

	connectionLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "loadbalancer_connection_latency_seconds",
			Help:    "Histogram of latency (seconds) for handling connections",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Go Runtime Metrics
var (
	goroutinesMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_goroutines_running",
			Help: "Number of currently running Go routines",
		},
	)

	gcPauseTotalMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "go_gc_pause_total_ms",
			Help: "Total time spent in garbage collection pauses (milliseconds)",
		},
	)

	memoryAllocMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_alloc_bytes",
			Help: "Current bytes of memory allocated by the Go runtime",
		},
	)

	goVersionMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_version_info",
			Help: "Go version info (static metric)",
		},
		[]string{"version"},
	)
)

func init() {
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(connectionLatency)

	prometheus.MustRegister(goroutinesMetric)
	prometheus.MustRegister(gcPauseTotalMetric)
	prometheus.MustRegister(memoryAllocMetric)
	prometheus.MustRegister(goVersionMetric)

	goVersionMetric.WithLabelValues(runtime.Version()).Set(1)

	go collectRuntimeMetrics()
}

func collectRuntimeMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var stats runtime.MemStats

	for range ticker.C {
		goroutinesMetric.Set(float64(runtime.NumGoroutine()))

		runtime.ReadMemStats(&stats)
		gcPauseTotalMetric.Add(float64(stats.PauseTotalNs / 1e6)) // Convert ns -> ms

		memoryAllocMetric.Set(float64(stats.Alloc))
	}
}
