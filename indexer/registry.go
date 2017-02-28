package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
	"sync"
	"time"
)

type Registry struct {
	Agents         map[string]models.Agent
	mutex          sync.RWMutex
	NatsConnection *nats.Conn
	Config         *Config
}

func NewRegistry(config *Config, nc *nats.Conn) *Registry {
	registry := &Registry{
		Agents:         make(map[string]models.Agent),
		NatsConnection: nc,
		Config:         config,
	}

	return registry
}

func (reg *Registry) WatchForAgents() {

	reg.NatsConnection.Subscribe(reg.Config.RegistryTopic, func(m *nats.Msg) {
		var payload models.Agent
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Agent payload unmarshal error %v", err)
		} else {
			log.Debugf("Agent payload received from %v", payload.Hostname)
			reg.mutex.Lock()
			defer reg.mutex.Unlock()
			reg.Agents[payload.Id] = payload
		}
	})
}

func (reg *Registry) GetActiveAgents() []models.Agent {
	agents := make([]models.Agent, 0)

	reg.mutex.Lock()
	defer reg.mutex.Unlock()

	for _, a := range reg.Agents {
		if a.LastSeen.After(time.Now().Add(-1 * time.Minute).UTC()) {
			agents = append(agents, a)
		}
	}

	return agents
}
