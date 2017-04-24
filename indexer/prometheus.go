package main

import "github.com/prometheus/client_golang/prometheus"

type Prometheus struct {
	requestsTotal   *prometheus.CounterVec
	requestsLatency *prometheus.SummaryVec
}

func NewPrometheus(namespace string, subsystem string) *Prometheus {
	prom := &Prometheus{}

	prom.requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "The number of requests.",
		},
		[]string{"method", "path", "status"},
	)

	prom.requestsLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "request_latency",
			Help:      "The latency of requests.",
		},
		[]string{"method", "path", "status"},
	)

	prometheus.MustRegister(prom.requestsTotal)
	prometheus.MustRegister(prom.requestsLatency)

	return prom
}
