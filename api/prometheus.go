package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stefanprodan/chi"
	"github.com/stefanprodan/chi/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			t2 := time.Now()
			// ignore websockets and the /metrics endpoint
			if !isWSRequest(r) && !isPrometheusRequest(r) {
				http_requests_total.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(ww.Status())).Inc()
				http_requests_latency.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(ww.Status())).Observe(t2.Sub(t1).Seconds())
			}
		}()

		next.ServeHTTP(ww, r)
	})
}

func PrometheusRegister() {
	prometheus.MustRegister(http_requests_total)
	prometheus.MustRegister(http_requests_latency)
}

func PrometheusMetrics(next http.Handler) http.Handler {
	return promhttp.Handler()
}

func PrometheusRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(PrometheusMetrics)
	r.Get("/", emptyHandler)
	return r
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {

}

var http_requests_total = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "syros",
		Subsystem: "api",
		Name:      "http_requests_total",
		Help:      "The number of HTTP requests.",
	},
	[]string{"method", "path", "status"},
)

var http_requests_latency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: "syros",
		Subsystem: "api",
		Name:      "http_request_latency",
		Help:      "The latency of HTTP requests.",
	},
	[]string{"method", "path", "status"},
)

func isWSRequest(req *http.Request) bool {
	return strings.ToLower(req.Header.Get("Upgrade")) == "websocket" &&
		strings.ToLower(req.Header.Get("Connection")) == "upgrade"
}

func isPrometheusRequest(req *http.Request) bool {
	return strings.Contains(strings.ToLower(req.URL.Path), "/metrics")
}
