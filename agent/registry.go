package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
	"os"
	"time"
)

type Registry struct {
	Topic          string
	Agent          models.SyrosService
	NatsConnection *nats.Conn
}

func NewRegistry(config *Config, nc *nats.Conn) *Registry {

	agent := models.SyrosService{
		Environment: config.Environment,
		Type:        "agent",
	}

	agent.Config, _ = models.ConfigToMap(config, "m")
	agent.Hostname, _ = os.Hostname()
	uuid, _ := models.NewUUID()
	agent.Id = models.Hash(agent.Hostname + uuid)

	registry := &Registry{
		Topic:          "registry",
		NatsConnection: nc,
		Agent:          agent,
	}

	return registry
}

func (r *Registry) RegisterAgent() error {

	r.Agent.Collected = time.Now().UTC()
	jsonPayload, err := json.Marshal(r.Agent)
	if err != nil {
		return fmt.Errorf("Agent payload marshal error %v", err)
	} else {
		err := r.NatsConnection.Publish(r.Topic, jsonPayload)
		if err != nil {
			return fmt.Errorf("Registry NATS publish failed %v", err)
		}
	}
	return nil
}
