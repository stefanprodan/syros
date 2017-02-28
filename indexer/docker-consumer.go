package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
)

type DockerConsumer struct {
	Config         *Config
	NatsConnection *nats.Conn
	Repository     *Repository
}

func NewDockerConsumer(config *Config, nc *nats.Conn, repo *Repository) (*DockerConsumer, error) {
	consumer := &DockerConsumer{
		Config:         config,
		NatsConnection: nc,
		Repository:     repo,
	}
	return consumer, nil
}

func (c *DockerConsumer) Consume() {
	c.NatsConnection.QueueSubscribe(c.Config.CollectorTopic, c.Config.CollectorQueue, func(m *nats.Msg) {
		var payload models.DockerPayload
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Docker payload unmarshal error %v", err)
		} else {
			log.Debugf("Docker payload received from host %v ID %v running containes %v", payload.Host.Name, payload.Host.Id, payload.Host.ContainersRunning)
			c.Repository.HostUpsert(payload.Host)

			for _, container := range payload.Containers {
				c.Repository.ContainerUpsert(container)
			}
		}
	})
}
