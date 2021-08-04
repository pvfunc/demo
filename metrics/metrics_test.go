package metrics

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusHandler(t *testing.T) {
	t.Parallel()

	assert.ObjectsAreEqual(PrometheusHandler(), promhttp.Handler())
}

func Test_makeMetricsResponseWriter(t *testing.T) {
	t.Parallel()

	var w http.ResponseWriter

	want := metricsResponseWriter{w, http.StatusOK}
	got := makeMetricsResponseWriter(w)

	assert.Equal(t, got, &want)
}

func TestWriteHeader(t *testing.T) { //nolint:paralleltest
	w := httptest.NewRecorder()
	mrw := &metricsResponseWriter{w, http.StatusOK}

	for i := 0; i < 550; i++ {
		i := i
		t.Run("Test code "+strconv.Itoa(i), func(t *testing.T) {
			mrw.WriteHeader(i)

			assert.Equal(t, mrw.statusCode, i)
		})
	}
}

func TestUrlToLabel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "test empty path",
			path: "",
			want: "/",
		},
		{
			name: "test root path",
			path: "/",
			want: "/",
		},
		{
			name: "test start-build path",
			path: "/start-build/",
			want: "/start-build",
		},
		{
			name: "test get-build-info path",
			path: "/get-build-info",
			want: "/get-build-info",
		},
		{
			name: "test health path",
			path: "/_/health",
			want: "/_/health",
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := urlToLabel(tt.path)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestMakeMiddlewareOptions(t *testing.T) {
	t.Parallel()

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

	serviceMetricOptions := &ServiceMetricOptions{
		Counter:   counter,
		Histogram: histogram,
	}

	want := MiddlewareOptions{
		ServiceMetrics: serviceMetricOptions,
	}

	got := MakeMiddlewareOptions()

	assert.IsType(t, got.ServiceMetrics.Counter, want.ServiceMetrics.Counter)
	assert.IsType(t, got.ServiceMetrics.Histogram, want.ServiceMetrics.Histogram)
}
