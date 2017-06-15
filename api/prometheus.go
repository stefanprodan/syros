package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type Prometheus struct {
	ignoreWebsockets    bool
	ignoreMetrics       bool
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestsLatency *prometheus.SummaryVec
}

func NewPrometheus(namespace string, subsystem string, ignoreWebsockets bool, ignoreMetrics bool) *Prometheus {
	prom := &Prometheus{
		ignoreMetrics:    ignoreMetrics,
		ignoreWebsockets: ignoreWebsockets,
	}

	prom.httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_requests_total",
			Help:      "The number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	prom.httpRequestsLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "http_request_latency",
			Help:      "The latency of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	prometheus.MustRegister(prom.httpRequestsTotal)
	prometheus.MustRegister(prom.httpRequestsLatency)

	return prom
}

func (p *Prometheus) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			t2 := time.Now()
			// ignore websockets and the /metrics endpoint
			if !p.isWSRequest(r) && !p.isPrometheusRequest(r) {
				p.httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(ww.Status())).Inc()
				p.httpRequestsLatency.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(ww.Status())).Observe(t2.Sub(t1).Seconds())
			}
		}()

		next.ServeHTTP(ww, r)
	})
}

func (p *Prometheus) Router() http.Handler {
	promHandler := func(next http.Handler) http.Handler { return promhttp.Handler() }
	emptyHandler := func(w http.ResponseWriter, r *http.Request) {}
	r := chi.NewRouter()
	r.Use(promHandler)
	r.Get("/", emptyHandler)
	return r
}

func (p *Prometheus) isWSRequest(req *http.Request) bool {
	if p.ignoreWebsockets {
		return strings.ToLower(req.Header.Get("Upgrade")) == "websocket" &&
			strings.ToLower(req.Header.Get("Connection")) == "upgrade"
	} else {
		return false
	}
}

func (p *Prometheus) isPrometheusRequest(req *http.Request) bool {
	if p.ignoreMetrics {
		return strings.Contains(strings.ToLower(req.URL.Path), "/metrics")
	} else {
		return false
	}
}
