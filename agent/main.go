package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
)

var version = "undefined"

func main() {
	var config = &Config{}
	flag.StringVar(&config.Environment, "Environment", "dev", "Environment dev|int|stg|test|prep|prod")
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8886, "HTTP port to listen on")
	flag.IntVar(&config.CollectInterval, "CollectInterval", 10, "Collect interval in seconds")
	flag.StringVar(&config.DockerApiAddresses, "DockerApiAddresses", "unix:///var/run/docker.sock", "Docker hosts API addresses comma delimited")
	flag.StringVar(&config.ConsulApiAddresses, "ConsulApiAddresses", "", "Consul hosts API addresses comma delimited")
	flag.StringVar(&config.VSphereApiAddress, "VSphereApiAddress", "", "VSphere API address")
	flag.StringVar(&config.VSphereInclude, "VSphereInclude", "", "VM include filter comma delimited")
	flag.StringVar(&config.VSphereExclude, "VSphereExclude", "", "VM exclude filter comma delimited")
	flag.IntVar(&config.VSphereCollectInterval, "VSphereCollectInterval", 120, "vSphere collect interval in seconds")
	flag.StringVar(&config.Nats, "Nats", "nats://localhost:4222", "Nats server addresses comma delimited")
	flag.Parse()

	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	if len(config.VSphereApiAddress) > 0 {
		vCol, err := NewVSphereCollector(config.VSphereApiAddress, config.VSphereInclude, config.VSphereExclude, config.Environment, config.VSphereCollectInterval)
		_, err = vCol.Collect()
		if err != nil {
			log.Errorf("VSphere error %v", err)
		}
	}

	//nc, err := NewNatsConnection(config.Nats)
	//defer nc.Close()
	//if err != nil {
	//	log.Fatalf("Nats connection error %v", err)
	//}
	//log.Infof("Connected to NATS server %v status %v", nc.ConnectedUrl(), nc.Status())

	cronJob := cron.New()

	registry := NewRegistry(config, nil, cronJob)
	log.Infof("Register service as %v", registry.Agent.Id)
	registry.Register()

	coordinator, err := NewCoordinator(config, nil, cronJob)
	if err != nil {
		log.Fatalf("Coordinator error %v", err)
	}
	coordinator.Register()

	server := &HttpServer{
		Config: config,
	}
	log.Infof("Starting HTTP server on port %v", config.Port)
	go server.Start()

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Infof("Shuting down %v signal received", sig)
	coordinator.Deregister()
	time.Sleep(1 * time.Second)
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
