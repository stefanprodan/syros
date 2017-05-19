package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
)

type Consumer struct {
	Config         *Config
	NatsConnection *nats.EncodedConn
	Repository     *Repository
	metrics        *Prometheus
	dockerChan     chan *models.DockerPayload
	consulChan     chan *models.ConsulPayload
	vsphereChan    chan *models.VSpherePayload
}

func NewConsumer(config *Config, nc *nats.EncodedConn, repo *Repository, buffer int) (*Consumer, error) {
	consumer := &Consumer{
		Config:         config,
		NatsConnection: nc,
		Repository:     repo,
		dockerChan:     make(chan *models.DockerPayload, buffer),
		consulChan:     make(chan *models.ConsulPayload, buffer),
		vsphereChan:    make(chan *models.VSpherePayload, buffer),
	}

	consumer.metrics = NewPrometheus("syros", "indexer")

	return consumer, nil
}

func (c *Consumer) Consume() {
	c.DockerConsume()
	c.ConsulConsume()
	c.VSphereConsume()
}

func (c *Consumer) DockerConsume() {
	c.NatsConnection.BindRecvQueueChan("docker", c.Config.CollectorQueue, c.dockerChan)
	go func() {
		for {
			select {
			case payload := <-c.dockerChan:
				dockerSave(payload, c)
			}
		}
	}()
}

func dockerSave(payload *models.DockerPayload, c *Consumer) {
	status := "200"
	t1 := time.Now()
	if payload == nil {
		log.Error("Docker payload is nil")
		status = "500"
	} else {
		log.Debugf("Docker payload received from host %v running containes %v", payload.Host.Name, payload.Host.ContainersRunning)
		c.Repository.HostUpsert(payload.Host)
		c.Repository.ContainersUpsert(payload.Containers)
	}
	t2 := time.Now()
	c.metrics.requestsTotal.WithLabelValues("docker", c.Config.CollectorQueue, status).Inc()
	c.metrics.requestsLatency.WithLabelValues("docker", c.Config.CollectorQueue, status).Observe(t2.Sub(t1).Seconds())
}

func (c *Consumer) ConsulConsume() {
	c.NatsConnection.BindRecvQueueChan("consul", c.Config.CollectorQueue, c.consulChan)
	go func() {
		for {
			select {
			case payload := <-c.consulChan:
				consulSave(payload, c)
			}
		}
	}()
}

func consulSave(payload *models.ConsulPayload, c *Consumer) {
	status := "200"
	t1 := time.Now()
	if payload == nil {
		log.Error("Consul payload is nil")
	} else {
		log.Debugf("Consul payload received %v checks", len(payload.HealthChecks))
		c.Repository.ChecksUpsert(payload.HealthChecks)
	}
	t2 := time.Now()
	c.metrics.requestsTotal.WithLabelValues("consul", c.Config.CollectorQueue, status).Inc()
	c.metrics.requestsLatency.WithLabelValues("consul", c.Config.CollectorQueue, status).Observe(t2.Sub(t1).Seconds())
}

func (c *Consumer) VSphereConsume() {
	c.NatsConnection.BindRecvQueueChan("vsphere", c.Config.CollectorQueue, c.vsphereChan)
	go func() {
		for {
			select {
			case payload := <-c.vsphereChan:
				vsphereSave(payload, c)
			}
		}
	}()
}

func vsphereSave(payload *models.VSpherePayload, c *Consumer) {
	status := "200"
	t1 := time.Now()
	if payload == nil {
		log.Errorf("VSphere payload is nil")
		status = "500"
	} else {
		log.Debugf("VSphere payload received %v vms %v hosts %v datastores",
			len(payload.VMs), len(payload.Hosts), len(payload.DataStores))
		c.Repository.VSphereDatastoresUpsert(payload.DataStores)
		c.Repository.VSphereHostsUpsert(payload.Hosts)
		c.Repository.VSphereVMsUpsert(payload.VMs)
	}
	t2 := time.Now()
	c.metrics.requestsTotal.WithLabelValues("vsphere", c.Config.CollectorQueue, status).Inc()
	c.metrics.requestsLatency.WithLabelValues("vsphere", c.Config.CollectorQueue, status).Observe(t2.Sub(t1).Seconds())
}
