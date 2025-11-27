package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	BlocksProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "indexer_blocks_processed_total",
		Help: "Total number of blocks processed by the indexer.",
	})
	IndexingLagSeconds = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "indexer_lag_seconds",
		Help: "Difference between chain head and last indexed block in seconds.",
	})
	APILatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "api_request_latency_seconds",
		Help:    "Latency of API endpoints.",
		Buckets: prometheus.DefBuckets,
	})
	WSConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ws_connections",
		Help: "Number of active websocket connections.",
	})
)

func init() {
	prometheus.MustRegister(BlocksProcessed, IndexingLagSeconds, APILatency, WSConnections)
}
