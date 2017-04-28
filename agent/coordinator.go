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
	NatsConnection   *nats.Conn
	Config           *Config
	metrics          *Prometheus
}

func NewCoordinator(config *Config, nc *nats.Conn) (*Coordinator, error) {

	ep := make([]string, 0)
	co := &Coordinator{
		NatsConnection: nc,
		Config:         config,
	}
	co.metrics = NewPrometheus("syros", "agent")

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

	return co, nil
}

func (cor *Coordinator) StartDockerCollectors() {
	log.Infof("Starting %v Docker collector(s)", len(cor.DockerCollectors))
	for _, c := range cor.DockerCollectors {
		time.Sleep(100 * time.Millisecond)
		ticker := time.NewTicker(time.Duration(cor.Config.CollectInterval) * time.Second)
		go func(collector *DockerCollector) {
			for {
				select {
				case <-collector.StopChan:
					log.Infof("Collector exited %v", collector.ApiAddress)
					return
				case <-ticker.C:
					status := "200"
					t1 := time.Now()

					payload, err := collector.Collect()
					if err != nil {
						log.Errorf("Docker collector %v error %v", collector.ApiAddress, err)
					} else {
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

					t2 := time.Now()
					cor.metrics.requestsTotal.WithLabelValues("docker", collector.ApiAddress, status).Inc()
					cor.metrics.requestsLatency.WithLabelValues("docker", collector.ApiAddress, status).Observe(t2.Sub(t1).Seconds())
				}
			}
		}(c)
	}
}

func (cor *Coordinator) StartConsulCollectors() {
	log.Infof("Starting %v Consul collector(s)", len(cor.ConsulCollectors))
	for _, c := range cor.ConsulCollectors {
		ticker := time.NewTicker(time.Duration(cor.Config.CollectInterval) * time.Second)
		go func(collector *ConsulCollector) {
			for {
				select {
				case <-collector.StopChan:
					log.Infof("Collector exited %v", collector.ApiAddress)
					return
				case <-ticker.C:
					status := "200"
					t1 := time.Now()

					payload, err := collector.Collect()
					if err != nil {
						log.Errorf("Consul collector %v error %v", collector.ApiAddress, err)
					} else {
						jsonPayload, err := json.Marshal(payload)
						if err != nil {
							log.Errorf("Consul collector %v payload marshal error %v", collector.ApiAddress, err)
						} else {
							err := cor.NatsConnection.Publish(collector.Topic, jsonPayload)
							if err != nil {
								log.Errorf("Consul collector %v NATS publish failed %v", collector.ApiAddress, err)
							}
						}
					}

					t2 := time.Now()
					cor.metrics.requestsTotal.WithLabelValues("consul", collector.ApiAddress, status).Inc()
					cor.metrics.requestsLatency.WithLabelValues("consul", collector.ApiAddress, status).Observe(t2.Sub(t1).Seconds())
				}
			}
		}(c)
	}
}

func (cor *Coordinator) StartCollectors() {
	cor.StartDockerCollectors()
	cor.StartConsulCollectors()
}

func (cor *Coordinator) StopCollectors() {
	log.Infof("Stopping %v Docker collector(s)", len(cor.DockerCollectors))
	for _, c := range cor.DockerCollectors {
		c.StopChan <- true
	}

	log.Infof("Stopping %v Consul collector(s)", len(cor.ConsulCollectors))
	for _, c := range cor.ConsulCollectors {
		c.StopChan <- true
	}
}
