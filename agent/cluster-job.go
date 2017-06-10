package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

type clusterJob struct {
	collector *ClusterCollector
	nats      *nats.EncodedConn
	metrics   *Prometheus
	config    *Config
}

func (j clusterJob) Run() {
	status := "200"
	t1 := time.Now()

	payload, err := j.collector.Collect()
	if err != nil {
		status = "500"
		log.Errorf("Cluster collector %v error %v", j.collector.ApiAddress, err)
	} else {
		err = j.nats.Publish(j.collector.Topic, payload)
		if err != nil {
			status = "500"
			log.Errorf("Cluster collector %v Nats natsPublish error %v", j.collector.ApiAddress, err)
		}
	}

	t2 := time.Now()
	j.metrics.requestsTotal.WithLabelValues("cluster", j.collector.ApiAddress, status).Inc()
	j.metrics.requestsLatency.WithLabelValues("cluster", j.collector.ApiAddress, status).Observe(t2.Sub(t1).Seconds())
}
