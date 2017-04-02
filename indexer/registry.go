package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
	"sync"
)

type Registry struct {
	mutex          sync.RWMutex
	NatsConnection *nats.Conn
	Config         *Config
	Repository     *MongoRepository
}

func NewRegistry(config *Config, nc *nats.Conn, repo *MongoRepository) *Registry {
	registry := &Registry{
		NatsConnection: nc,
		Config:         config,
		Repository:     repo,
	}

	return registry
}

func (reg *Registry) WatchForAgents() {

	reg.NatsConnection.QueueSubscribe(reg.Config.RegistryTopic, reg.Config.RegistryQueue, func(m *nats.Msg) {
		var payload models.SyrosService
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("Agent payload unmarshal error %v", err)
		} else {
			log.Debugf("Agent payload received from %v", payload.Hostname)
			reg.Repository.SyrosServiceUpsert(payload)
		}
	})
}

func (reg *Registry) SelfRegister(indexer models.SyrosService) {
	reg.Repository.SyrosServiceUpsert(indexer)
}
