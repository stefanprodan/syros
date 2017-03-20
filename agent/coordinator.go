package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats"
	"strings"
	"time"
)

type Coordinator struct {
	DockerCollectors []*DockerCollector
	ConsulCollectors []*ConsulCollector
	Status           *AgentStatus
	NatsConnection   *nats.Conn
	Config           *Config
}

func NewCoordinator(config *Config, nc *nats.Conn) (*Coordinator, error) {

	ep := make([]string, 0)
	co := &Coordinator{
		NatsConnection: nc,
		Config:         config,
	}
	if len(config.DockerApiAddresses) > 0 {
		dh := strings.Split(config.DockerApiAddresses, ",")
		dc := make([]*DockerCollector, len(dh))
		for i, host := range dh {
			c, err := NewDockerCollector(host, config.Environment)
			if err != nil {
				return nil, err
			}
			dc[i] = c
		}
		ep = append(ep, dh...)
		co.DockerCollectors = dc
	}

	if len(config.ConsulApiAddresses) > 0 {
		ch := strings.Split(config.ConsulApiAddresses, ",")
		cc := make([]*ConsulCollector, len(ch))
		for i, host := range ch {
			c, err := NewConsulCollector(host, config.Environment)
			if err != nil {
				return nil, err
			}
			cc[i] = c
		}
		ep = append(ep, ch...)
		co.ConsulCollectors = cc
	}

	status, err := NewAgentStatus(ep)
	if err != nil {
		return nil, err
	}
	co.Status = status

	return co, nil
}

func (cor *Coordinator) StartDockerCollectors() {
	log.Infof("Starting %v Docker collector(s)", len(cor.DockerCollectors))
	for _, c := range cor.DockerCollectors {
		go func(collector *DockerCollector) {
			stop := false
			for !stop {
				select {
				case <-collector.StopChan:
					stop = true
				default:
					payload, err := collector.Collect()
					if err != nil {
						log.Errorf("Docker collector %v error %v", collector.ApiAddress, err)
						cor.Status.SetCollectorStatus(collector.ApiAddress, false)
					} else {
						cor.Status.SetCollectorStatus(collector.ApiAddress, true)
						jsonPayload, err := json.Marshal(payload)
						if err != nil {
							log.Errorf("Docker collector %v payload marshal error %v", collector.ApiAddress, err)
						} else {
							err := cor.NatsConnection.Publish(collector.Topic, jsonPayload)
							if err != nil {
								log.Errorf("Docker collector %v NATS publish failed %v", collector.ApiAddress, err)
							}
						}
					}
					time.Sleep(time.Duration(cor.Config.CollectInterval) * time.Second)
				}
			}
			log.Infof("Collector exited %v", collector.ApiAddress)
		}(c)
	}
}

func (cor *Coordinator) StartCollectors() {
	cor.StartDockerCollectors()
}

func (cor *Coordinator) StopCollectors() {
	log.Infof("Stopping %v Docker collector(s)", len(cor.DockerCollectors))
	for _, c := range cor.DockerCollectors {
		c.StopChan <- struct{}{}
	}

	log.Infof("Stopping %v Consul collector(s)", len(cor.ConsulCollectors))
	for _, c := range cor.ConsulCollectors {
		c.StopChan <- struct{}{}
	}
}
