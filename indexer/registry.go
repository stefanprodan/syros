package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/stefanprodan/syros/models"
	"os"
	"sync"
	"time"
)

type Registry struct {
	mutex          sync.RWMutex
	NatsConnection *nats.Conn
	Config         *Config
	Repository     *Repository
}

func NewRegistry(config *Config, nc *nats.Conn, repo *Repository) *Registry {
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

func (reg *Registry) Start() chan bool {

	indexer := models.SyrosService{
		Environment: "all",
		Type:        "indexer",
	}
	indexer.Config, _ = models.ConfigToMap(reg.Config, "m")
	indexer.Hostname, _ = os.Hostname()
	uuid, _ := models.NewUUID()
	indexer.Id = models.Hash(indexer.Hostname + uuid)

	stopped := make(chan bool, 1)
	ticker := time.NewTicker(10 * time.Second)

	go func(i models.SyrosService) {
		for {
			select {
			case <-ticker.C:
				i.Collected = time.Now().UTC()
				reg.Repository.SyrosServiceUpsert(i)
			case <-stopped:
				return
			}
		}
	}(indexer)

	return stopped
}
