package metric

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	responseTime    *prometheus.HistogramVec
	totalRequests   *prometheus.CounterVec
	duration        *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	handlerStatuses *prometheus.CounterVec
}

// NewMetrics creates new custom Prometheus metrics
func NewMetricsBuilder(app, metricsPrefix, version, hash, date string) *Metrics {
	labels := map[string]string{
		"app":       app,
		"version":   version,
		"hash":      hash,
		"buildTime": date,
	}

	pm := &Metrics{
		responseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        metricsPrefix + "response_time_seconds",
				Help:        "Description",
				ConstLabels: labels,
			},
			[]string{"endpoint"},
		),
		totalRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        metricsPrefix + "requests_total",
			Help:        "number of requests",
			ConstLabels: labels,
		}, []string{"code", "method"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        metricsPrefix + "requests_duration_seconds",
			Help:        "duration of a requests in seconds",
			ConstLabels: labels,
		}, []string{"code", "method"}),
		responseSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        metricsPrefix + "response_size_bytes",
			Help:        "size of the responses in bytes",
			ConstLabels: labels,
		}, []string{"code", "method"}),
		requestSize: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        metricsPrefix + "requests_size_bytes",
			Help:        "size of the requests in bytes",
			ConstLabels: labels,
		}, []string{"code", "method"}),
		handlerStatuses: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        metricsPrefix + "requests_statuses_total",
			Help:        "count number of responses per status",
			ConstLabels: labels,
		}, []string{"method", "status_bucket"}),
	}

	err := prometheus.Register(pm)
	if e := new(prometheus.AlreadyRegisteredError); errors.As(err, e) {
		return pm
	} else if err != nil {
		panic(err)
	}

	return pm
}

// Describe implements prometheus Collector interface.
func (h *Metrics) Describe(in chan<- *prometheus.Desc) {
	h.duration.Describe(in)
	h.totalRequests.Describe(in)
	h.requestSize.Describe(in)
	h.responseSize.Describe(in)
	h.handlerStatuses.Describe(in)
	h.responseTime.Describe(in)
}

// Collect implements prometheus Collector interface.
func (h *Metrics) Collect(in chan<- prometheus.Metric) {
	h.duration.Collect(in)
	h.totalRequests.Collect(in)
	h.requestSize.Collect(in)
	h.responseSize.Collect(in)
	h.handlerStatuses.Collect(in)
	h.responseTime.Collect(in)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (h Metrics) instrumentHandlerStatusBucket(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// res := 200
		// res := &loggingResponseWriter{rw, 0}
		next.ServeHTTP(rw, r)

		res := rw.(*loggingResponseWriter)
		fmt.Println(res)

		statusBucket := "unknown"
		switch {
		case res.statusCode >= 200 && res.statusCode <= 299:
			statusBucket = "2xx"
		case res.statusCode >= 300 && res.statusCode <= 399:
			statusBucket = "3xx"
		case res.statusCode >= 400 && res.statusCode <= 499:
			statusBucket = "4xx"
		case res.statusCode >= 500 && res.statusCode <= 599:
			statusBucket = "5xx"
		}

		h.handlerStatuses.With(prometheus.Labels{"method": r.Method, "status_bucket": statusBucket}).
			Inc()
	}
}

// Instrument will instrument any http.HandlerFunc with custom metrics
func (h Metrics) Instrument(rw http.ResponseWriter, next http.HandlerFunc, endpoint string) http.HandlerFunc {
	labels := prometheus.Labels{}

	res, ok := rw.(*loggingResponseWriter)
	if ok && res.statusCode != 0 {
		fmt.Println("http status code :", res.statusCode)
		labels = prometheus.Labels{"code": strconv.Itoa(res.statusCode)}
	}
	fmt.Println("http status code :", res.statusCode)

	wrapped := promhttp.InstrumentHandlerResponseSize(h.responseSize.MustCurryWith(labels), next)
	wrapped = promhttp.InstrumentHandlerCounter(h.totalRequests.MustCurryWith(labels), wrapped)
	wrapped = promhttp.InstrumentHandlerDuration(h.duration.MustCurryWith(labels), wrapped)
	wrapped = promhttp.InstrumentHandlerDuration(h.responseTime.MustCurryWith(prometheus.Labels{"endpoint": endpoint}), wrapped)
	wrapped = promhttp.InstrumentHandlerRequestSize(h.requestSize.MustCurryWith(labels), wrapped)
	wrapped = h.instrumentHandlerStatusBucket(wrapped)

	return wrapped.ServeHTTP
}
