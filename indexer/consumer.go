package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
	"time"
)

type Consumer struct {
	Config         *Config
	NatsConnection *nats.Conn
	Repository     *Repository
	metrics        *Prometheus
}

func NewConsumer(config *Config, nc *nats.Conn, repo *Repository) (*Consumer, error) {
	consumer := &Consumer{
		Config:         config,
		NatsConnection: nc,
		Repository:     repo,
	}

	consumer.metrics = NewPrometheus("syros", "indexer")

	return consumer, nil
}

func (c *Consumer) Consume() {
	c.DockerConsume()
	c.ConsulConsume()
}

func (c *Consumer) DockerConsume() {
	c.NatsConnection.QueueSubscribe("docker", c.Config.CollectorQueue, func(m *nats.Msg) {
		status := "200"
		t1 := time.Now()

		var payload models.DockerPayload
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Docker payload unmarshal error %v", err)
			status = "500"
		} else {
			log.Debugf("Docker payload received from host %v running containes %v", payload.Host.Name, payload.Host.ContainersRunning)
			c.Repository.HostUpsert(payload.Host)
			c.Repository.ContainersUpsert(payload.Containers)
		}

		t2 := time.Now()
		c.metrics.requestsTotal.WithLabelValues("docker", c.Config.CollectorQueue, status).Inc()
		c.metrics.requestsLatency.WithLabelValues("docker", c.Config.CollectorQueue, status).Observe(t2.Sub(t1).Seconds())
	})
}

func (c *Consumer) ConsulConsume() {
	c.NatsConnection.QueueSubscribe("consul", c.Config.CollectorQueue, func(m *nats.Msg) {
		status := "200"
		t1 := time.Now()

		var payload models.ConsulPayload
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Consul payload unmarshal error %v", err)
		} else {
			log.Debugf("Consul payload received %v checks", len(payload.HealthChecks))
			c.Repository.ChecksUpsert(payload.HealthChecks)
		}

		t2 := time.Now()
		c.metrics.requestsTotal.WithLabelValues("consul", c.Config.CollectorQueue, status).Inc()
		c.metrics.requestsLatency.WithLabelValues("consul", c.Config.CollectorQueue, status).Observe(t2.Sub(t1).Seconds())
	})
}
