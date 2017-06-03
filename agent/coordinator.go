package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"github.com/robfig/cron"
)

type Coordinator struct {
	NatsConnection  *nats.EncodedConn
	Config          *Config
	CollectorConfig *CollectorConfig
	Cron            *cron.Cron
	metrics         *Prometheus
}

func NewCoordinator(config *Config, collector *CollectorConfig, nc *nats.EncodedConn, cron *cron.Cron) (*Coordinator, error) {
	co := &Coordinator{
		NatsConnection:  nc,
		Cron:            cron,
		Config:          config,
		CollectorConfig: collector,
	}
	co.metrics = NewPrometheus("syros", "agent")

	return co, nil
}

func (cor *Coordinator) Register() {

	for _, c := range cor.CollectorConfig.Docker.Endpoints {
		col, err := NewDockerCollector(c, cor.Config.Environment)
		if err != nil {
			log.Errorf("Collector %v init error", c)
		} else {
			cor.Cron.AddJob(cor.CollectorConfig.Docker.Cron,
				dockerJob{col, cor.NatsConnection, cor.metrics, cor.Config})
		}
	}

	for _, c := range cor.CollectorConfig.Consul.Endpoints {
		col, err := NewConsulCollector(c, cor.Config.Environment)
		if err != nil {
			log.Errorf("Collector %v init error", c)
		} else {
			cor.Cron.AddJob(cor.CollectorConfig.Consul.Cron,
				consulJob{col, cor.NatsConnection, cor.metrics, cor.Config})
		}
	}

	for _, c := range cor.CollectorConfig.VSphere.Endpoints {
		col, err := NewVSphereCollector(c, cor.CollectorConfig.VSphere.Include, cor.CollectorConfig.VSphere.Exclude, cor.Config.Environment)
		if err != nil {
			log.Errorf("Collector %v init error", c)
		} else {
			cor.Cron.AddJob(cor.CollectorConfig.VSphere.Cron,
				vsphereJob{col, cor.NatsConnection, cor.metrics, cor.Config})
		}
	}

	cor.Cron.Start()
}

func (cor *Coordinator) Deregister() {
	cor.Cron.Stop()
}
