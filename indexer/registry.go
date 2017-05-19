package main

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
)

type Registry struct {
	NatsConnection *nats.EncodedConn
	Config         *Config
	Repository     *Repository
}

func NewRegistry(config *Config, nc *nats.EncodedConn, repo *Repository) *Registry {
	registry := &Registry{
		NatsConnection: nc,
		Config:         config,
		Repository:     repo,
	}

	return registry
}

func (reg *Registry) WatchForAgents() {

	reg.NatsConnection.QueueSubscribe(reg.Config.RegistryTopic, reg.Config.RegistryQueue, func(payload *models.SyrosService) {
		if payload == nil {
			log.Error("Agent payload is nil")
		} else {
			log.Debugf("Agent payload received from %v", payload.Hostname)
			go reg.Repository.SyrosServiceUpsert(*payload)
		}
	})
}

func (reg *Registry) Start() chan bool {

	indexer := models.SyrosService{
		Environment: "all",
		Type:        "indexer",
	}
	indexer.Config, _ = models.ConfigToMap(reg.Config, "m")
	indexer.Config["syros_version"] = version
	indexer.Config["os"] = runtime.GOOS
	indexer.Config["arch"] = runtime.GOARCH
	indexer.Config["golang"] = runtime.Version()
	indexer.Config["max_procs"] = strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10)
	indexer.Config["goroutines"] = strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
	indexer.Config["cpu_count"] = strconv.FormatInt(int64(runtime.NumCPU()), 10)
	indexer.Hostname, _ = os.Hostname()
	uuid, _ := models.NewUUID()
	indexer.Id = models.Hash(indexer.Hostname + uuid)

	stopped := make(chan bool, 1)
	ticker := time.NewTicker(15 * time.Second)

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
