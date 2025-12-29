package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultRegistry pre-registers standard collectors.
var DefaultRegistry = prometheus.NewRegistry()

func init() {
	DefaultRegistry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	DefaultRegistry.MustRegister(prometheus.NewGoCollector())
}

func Handler() http.Handler {
	return promhttp.HandlerFor(DefaultRegistry, promhttp.HandlerOpts{})
}
