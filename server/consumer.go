package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/nats-io/go-nats"
)

type DockerConsumer struct {
	Config         *Config
	NatsConnection *nats.Conn
}

type DockerPayload struct {
	Host              types.Info
	ContainerList     []types.Container
	ContainerInfoList []types.ContainerJSON
}

func NewNatsConnection(servers string) (*nats.Conn, error) {
	nc, err := nats.Connect(servers,
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Warnf("Got disconnected from NATS %v", servers)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("Got reconnected to NATS %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Errorf("NATS connection closed. Reason: %q", nc.LastError())
		}),
	)
	return nc, err
}

func NewDockerConsumer(config *Config, nc *nats.Conn) (*DockerConsumer, error) {
	consumer := &DockerConsumer{
		Config:         config,
		NatsConnection: nc,
	}
	return consumer, nil
}

func (c *DockerConsumer) Consume() {
	c.NatsConnection.QueueSubscribe(c.Config.CollectorTopic, c.Config.CollectorQueue, func(m *nats.Msg) {
		var payload DockerPayload
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("DockerPayload unmarshal error %v", err)
		} else {
			log.Infof("Host %v running containes %v", payload.Host.Name, payload.Host.ContainersRunning)
		}
	})
}
