package security

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Middleware interface {
	WrapHandler(handlerName string, handler http.Handler) http.HandlerFunc
}

type middleware struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.SummaryVec
	responseSize    *prometheus.SummaryVec
}

func (m *middleware) WrapHandler(handlerName string, handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

		handler.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		method := r.Method
		code := strconv.Itoa(wrapped.statusCode)

		m.requestsTotal.WithLabelValues(handlerName, method, code).Inc()
		m.requestDuration.WithLabelValues(handlerName, method, code).Observe(duration)

		if r.ContentLength > 0 {
			m.requestSize.WithLabelValues(handlerName, method, code).Observe(float64(r.ContentLength))
		}
		if wrapped.size > 0 {
			m.responseSize.WithLabelValues(handlerName, method, code).Observe(float64(wrapped.size))
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

func New(registry prometheus.Registerer, buckets []float64) Middleware {
	if buckets == nil {
		buckets = prometheus.ExponentialBuckets(0.1, 1.5, 5)
	}

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method", "code"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: buckets,
		},
		[]string{"handler", "method", "code"},
	)

	requestSize := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Size of HTTP requests in bytes",
		},
		[]string{"handler", "method", "code"},
	)

	responseSize := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Size of HTTP responses in bytes",
		},
		[]string{"handler", "method", "code"},
	)

	registry.MustRegister(requestsTotal, requestDuration, requestSize, responseSize)

	return &middleware{
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
		requestSize:     requestSize,
		responseSize:    responseSize,
	}
}
