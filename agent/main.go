package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var config = &Config{}
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8000, "HTTP port to listen on")
	flag.IntVar(&config.CollectInterval, "CollectInterval", 10, "Collect interval in seconds")
	flag.StringVar(&config.DockerApiAddresses, "DockerApiAddresses", "", "Docker hosts API addresses comma delimited")
	flag.StringVar(&config.Nats, "Nats", "nats://localhost:4222", "Nats server addresses comma delimited")
	flag.StringVar(&config.CollectorTopic, "CollectorTopic", "docker", "Nats collector topic name")
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
		collector, err := NewDockerCollector(host)
		if err != nil {
			log.Fatal(err)
		}
		collectors[i] = collector
	}

	status, err := NewAgentStatus(hosts)
	if err != nil {
		log.Fatal(err)
	}

	agent := models.Agent{}
	agent.Hostname, _ = os.Hostname()
	agent.Id, _ = newUUID()
	log.Infof("Register service as %v", agent.Hostname)
	go func(a models.Agent) {
		for true {
			agent.LastSeen = time.Now().UTC()
			jsonPayload, err := json.Marshal(agent)
			if err != nil {
				log.Errorf("Agent payload marshal error %v", err)
			} else {
				err := nc.Publish(config.RegistryTopic, jsonPayload)
				if err != nil {
					log.Errorf("Registry NATS publish failed %v", err)
				}
			}
			time.Sleep(5 * time.Second)
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
						log.Errorf("Docker collector % error %v", collector.ApiAddress, err)
						status.SetCollectorStatus(collector.ApiAddress, false, nil)
					} else {
						status.SetCollectorStatus(collector.ApiAddress, true, payload)
						jsonPayload, err := json.Marshal(payload)
						if err != nil {
							log.Errorf("Docker collector %v payload marshal error %v", collector.ApiAddress, err)
						} else {
							err := nc.Publish(config.CollectorTopic, jsonPayload)
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

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
