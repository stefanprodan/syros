package main

import (
	"encoding/json"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var config = &Config{}
	flag.StringVar(&config.Environment, "Environment", "dev", "Environment dev|int|stg|test|prep|prod")
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8886, "HTTP port to listen on")
	flag.IntVar(&config.CollectInterval, "CollectInterval", 10, "Collect interval in seconds")
	flag.StringVar(&config.DockerApiAddresses, "DockerApiAddresses", "unix:///var/run/docker.sock", "Docker hosts API addresses comma delimited")
	flag.StringVar(&config.Nats, "Nats", "nats://localhost:4222", "Nats server addresses comma delimited")
	flag.StringVar(&config.RegistryTopic, "RegistryTopic", "registry", "Nats registry topic name")
	flag.Parse()

	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	nc, err := NewNatsConnection(config.Nats)
	defer nc.Close()

	if err != nil {
		log.Fatalf("Nats connection error %v", err)
	}
	log.Infof("Connected to NATS server %v status %v", nc.ConnectedUrl(), nc.Status())

	hosts := strings.Split(config.DockerApiAddresses, ",")
	if len(hosts) < 1 {
		log.Fatalf("no hosts supplied %s", config.DockerApiAddresses)
	}

	collectors := make([]*DockerCollector, len(hosts))
	for i, host := range hosts {
		collector, err := NewDockerCollector(host, config.Environment)
		if err != nil {
			log.Fatal(err)
		}
		collectors[i] = collector
	}

	status, err := NewAgentStatus(hosts)
	if err != nil {
		log.Fatal(err)
	}

	agent := models.SyrosService{
		Environment: config.Environment,
		Type:        "agent",
	}
	agent.Config, _ = models.ConfigToMap(config, "m")
	agent.Hostname, _ = os.Hostname()
	uuid, _ := models.NewUUID()
	agent.Id = models.Hash(agent.Hostname + uuid)
	log.Infof("Register service as %v", agent.Hostname)
	go func(a models.SyrosService) {
		for true {
			agent.Collected = time.Now().UTC()
			jsonPayload, err := json.Marshal(agent)
			if err != nil {
				log.Errorf("Agent payload marshal error %v", err)
			} else {
				err := nc.Publish(config.RegistryTopic, jsonPayload)
				if err != nil {
					log.Errorf("Registry NATS publish failed %v", err)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}(agent)

	log.Infof("Starting %v collector(s), collect interval is set to %v second(s)", len(collectors), config.CollectInterval)
	for _, c := range collectors {
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
						status.SetCollectorStatus(collector.ApiAddress, false, nil)
					} else {
						status.SetCollectorStatus(collector.ApiAddress, true, payload)
						jsonPayload, err := json.Marshal(payload)
						if err != nil {
							log.Errorf("Docker collector %v payload marshal error %v", collector.ApiAddress, err)
						} else {
							err := nc.Publish(collector.Topic, jsonPayload)
							if err != nil {
								log.Errorf("Docker collector %v NATS publish failed %v", collector.ApiAddress, err)
							}
						}
					}
					time.Sleep(time.Duration(config.CollectInterval) * time.Second)
				}
			}
			log.Infof("Collector exited %v", collector.ApiAddress)
		}(c)
	}

	server := &HttpServer{
		Config: config,
		Status: status,
	}
	log.Infof("Starting HTTP server on port %v", config.Port)
	go server.Start()

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Infof("Shuting down %v signal received", sig)
	server.Stop()
	log.Infof("Stopping %v collector(s)", len(collectors))
	for _, collector := range collectors {
		collector.StopChan <- struct{}{}
	}
	time.Sleep(10 * time.Second)
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
