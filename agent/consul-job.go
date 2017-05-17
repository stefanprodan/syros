package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

type consulJob struct {
	collector *ConsulCollector
	nats      *nats.Conn
	metrics   *Prometheus
	config    *Config
}

func (j consulJob) Run() {
	status := "200"
	t1 := time.Now()

	payload, err := j.collector.Collect()
	if err != nil {
		status = "500"
		log.Errorf("Consul collector %v error %v", j.collector.ApiAddress, err)
	} else {
		err = natsPublish(j.config.Nats, j.collector.Topic, payload)
		if err != nil {
			status = "500"
			log.Errorf("Consul collector %v Nats natsPublish error %v", j.collector.ApiAddress, err)
		}
	}

	t2 := time.Now()
	j.metrics.requestsTotal.WithLabelValues("consul", j.collector.ApiAddress, status).Inc()
	j.metrics.requestsLatency.WithLabelValues("consul", j.collector.ApiAddress, status).Observe(t2.Sub(t1).Seconds())
}
