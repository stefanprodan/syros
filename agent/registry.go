package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/robfig/cron"
	"github.com/stefanprodan/syros/models"
	"os"
	"runtime"
	"strconv"
	"time"
)

type Registry struct {
	Topic          string
	Agent          models.SyrosService
	NatsConnection *nats.Conn
	Cron           *cron.Cron
	Config           *Config
}

func NewRegistry(config *Config, nc *nats.Conn, cron *cron.Cron) *Registry {

	agent := models.SyrosService{
		Environment: config.Environment,
		Type:        "agent",
	}

	agent.Config, _ = models.ConfigToMap(config, "m")
	agent.Config["syros_version"] = version
	agent.Config["os"] = runtime.GOOS
	agent.Config["arch"] = runtime.GOARCH
	agent.Config["golang"] = runtime.Version()
	agent.Config["max_procs"] = strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10)
	agent.Config["goroutines"] = strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
	agent.Config["cpu_count"] = strconv.FormatInt(int64(runtime.NumCPU()), 10)
	agent.Hostname, _ = os.Hostname()
	uuid, _ := models.NewUUID()
	agent.Id = models.Hash(agent.Hostname + uuid)

	registry := &Registry{
		Topic:          "registry",
		NatsConnection: nc,
		Agent:          agent,
		Cron:           cron,
		Config: config,
	}

	return registry
}

func (r *Registry) Register() {
	r.Cron.AddFunc("10 * * * *", func() {
		err := r.RegisterAgent()
		if err != nil {
			log.Error(err)
		}
	})
}

func (r *Registry) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := r.RegisterAgent()
				if err != nil {
					log.Error(err)
				}
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}

func (r *Registry) RegisterAgent() error {
	ag := r.Agent
	ag.Collected = time.Now().UTC()

	jsonPayload, err := json.Marshal(ag)
	if err != nil {
		return fmt.Errorf("Agent payload marshal error %v", err)
	} else {
		nc, err := nats.Connect(r.Config.Nats)
		if err != nil {
			log.Errorf("Registry NATS publish failed %v", err)
		}
		err = nc.Publish(r.Topic, jsonPayload)
		if err != nil {
			return fmt.Errorf("Registry NATS publish failed %v", err)
		}
		if nc != nil && !nc.IsClosed() {
			nc.Close()
		}
	}
	return nil
}
