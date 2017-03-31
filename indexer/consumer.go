package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
)

type Consumer struct {
	Config         *Config
	NatsConnection *nats.Conn
	Repository     *Repository
}

func NewConsumer(config *Config, nc *nats.Conn, repo *Repository) (*Consumer, error) {
	consumer := &Consumer{
		Config:         config,
		NatsConnection: nc,
		Repository:     repo,
	}
	return consumer, nil
}

func (c *Consumer) Consume() {
	c.DockerConsume()
	c.ConsulConsume()
}

func (c *Consumer) DockerConsume() {
	c.NatsConnection.QueueSubscribe("docker", c.Config.CollectorQueue, func(m *nats.Msg) {
		crepo, err := NewRepository(c.Config)
		if err != nil {
			log.Fatalf("RethinkDB connection error %v", err)
		}
		var payload models.DockerPayload
		err = json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Docker payload unmarshal error %v", err)
		} else {
			log.Debugf("Docker payload received from host %v running containes %v", payload.Host.Name, payload.Host.ContainersRunning)
			crepo.HostUpsert(payload.Host)

			for _, container := range payload.Containers {
				crepo.ContainerUpsert(container)
			}
			crepo.Session.Close()
		}
	})
}

func (c *Consumer) ConsulConsume() {
	c.NatsConnection.QueueSubscribe("consul", c.Config.CollectorQueue, func(m *nats.Msg) {
		crepo, err := NewRepository(c.Config)
		if err != nil {
			log.Fatalf("RethinkDB connection error %v", err)
		}

		var payload models.ConsulPayload
		err = json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Consul payload unmarshal error %v", err)
		} else {
			log.Debugf("Consul payload received %v checks", len(payload.HealthChecks))
			for _, check := range payload.HealthChecks {
				crepo.CheckUpsert(check)
			}
			crepo.Session.Close()
		}
	})
}
