package metrics

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusHandler Bootstraps prometheus for metrics collection
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}

// MiddlewareOptions metrics
type MiddlewareOptions struct {
	ServiceMetrics *ServiceMetricOptions
}

// ServiceMetricOptions provides RED metrics
type ServiceMetricOptions struct {
	Histogram *prometheus.HistogramVec
	Counter   *prometheus.CounterVec
}

func makeMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{w, http.StatusOK}
}

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

// MiddlewareHandler middleware handler for router
func (mmw *MiddlewareOptions) MiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mrw := makeMetricsResponseWriter(w)
		then := time.Now()
		next.ServeHTTP(mrw, r)

		code := fmt.Sprintf("%d", mrw.statusCode)
		path := urlToLabel(r.URL.Path)

		mmw.ServiceMetrics.Counter.WithLabelValues(r.Method, path, code).Inc()
		mmw.ServiceMetrics.Histogram.WithLabelValues(r.Method, path, code).Observe(time.Since(then).Seconds())
	})
}

func urlToLabel(path string) string {
	if len(path) > 0 {
		path = strings.TrimRight(path, "/")
	}

	if path == "" {
		path = "/"
	}

	return path
}

// Synchronize to make sure MustRegister only called once
var once = sync.Once{}

// MakeMiddlewareOptions make metrics for service
func MakeMiddlewareOptions() MiddlewareOptions {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Seconds spent serving HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status"})

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "The total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	once.Do(func() {
		prometheus.MustRegister(histogram)
		prometheus.MustRegister(counter)
	})

	serviceMetricOptions := &ServiceMetricOptions{
		Counter:   counter,
		Histogram: histogram,
	}

	middlewareOptions := MiddlewareOptions{
		ServiceMetrics: serviceMetricOptions,
	}

	return middlewareOptions
}
