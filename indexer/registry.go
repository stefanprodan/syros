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
	Agents         map[string]models.SyrosService
	mutex          sync.RWMutex
	NatsConnection *nats.Conn
	Config         *Config
	Repository     *Repository
}

func NewRegistry(config *Config, nc *nats.Conn, repo *Repository) *Registry {
	registry := &Registry{
		Agents:         make(map[string]models.SyrosService),
		NatsConnection: nc,
		Config:         config,
		Repository:     repo,
	}

	return registry
}

func (reg *Registry) WatchForAgents() {

	reg.NatsConnection.Subscribe(reg.Config.RegistryTopic, func(m *nats.Msg) {
		var payload models.SyrosService
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Agent payload unmarshal error %v", err)
		} else {
			log.Debugf("Agent payload received from %v", payload.Hostname)
			reg.Repository.SyrosServiceUpsert(payload)
			reg.mutex.Lock()
			defer reg.mutex.Unlock()
			reg.Agents[payload.Id] = payload
		}
	})
}

func (reg *Registry) SelfRegister(indexer models.SyrosService) {
	reg.Repository.SyrosServiceUpsert(indexer)
}

func (reg *Registry) GetActiveAgents() []models.SyrosService {
	agents := make([]models.SyrosService, 0)

	reg.mutex.Lock()
	defer reg.mutex.Unlock()

	for _, a := range reg.Agents {
		if a.Collected.After(time.Now().Add(-1 * time.Minute).UTC()) {
			agents = append(agents, a)
		}
	}

	return agents
}
