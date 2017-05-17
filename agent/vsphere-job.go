package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

type vsphereJob struct {
	collector *VSphereCollector
	nats      *nats.Conn
	metrics   *Prometheus
	config    *Config
}

func (j vsphereJob) Run() {
	status := "200"
	t1 := time.Now()

	payload, err := j.collector.Collect()
	if err != nil {
		status = "500"
		log.Errorf("vSphere collector %v error %v", j.collector.ApiAddress, err)
	} else {
		err = natsPublish(j.config.Nats, j.collector.Topic, payload)
		if err != nil {
			status = "500"
			log.Errorf("vSphere collector %v Nats natsPublish error %v", j.collector.ApiAddress, err)
		}
	}

	t2 := time.Now()
	j.metrics.requestsTotal.WithLabelValues("vsphere", j.collector.ApiAddress, status).Inc()
	j.metrics.requestsLatency.WithLabelValues("vsphere", j.collector.ApiAddress, status).Observe(t2.Sub(t1).Seconds())
}
