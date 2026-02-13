package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EventsIngested = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "edtech_events_ingested_total",
			Help: "Total number of events ingested",
		},
		[]string{"type", "status"},
	)

	WorkerProcessingLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "edtech_worker_processing_seconds",
			Help:    "Worker event processing latency",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
		},
		[]string{"type"},
	)

	WorkerFailures = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "edtech_worker_failures_total",
			Help: "Total worker processing failures",
		},
		[]string{"type"},
	)

	DashboardQueryLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "edtech_dashboard_query_seconds",
			Help:    "Dashboard API query latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	IngestionLatency = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "edtech_ingestion_seconds",
			Help:    "Event ingestion request latency",
			Buckets: prometheus.DefBuckets,
		},
	)
)
