package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var version = "undefined"

func main() {
	var config = &Config{}
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8887, "HTTP port to listen on")
	flag.StringVar(&config.Nats, "Nats", "nats://localhost:4222", "Nats server addresses comma delimited")
	flag.StringVar(&config.CollectorQueue, "CollectorQueue", "syros", "Nats collector queue name")
	flag.StringVar(&config.RegistryTopic, "RegistryTopic", "registry", "Nats registry topic name")
	flag.StringVar(&config.RegistryQueue, "RegistryQueue", "syros", "Nats registry queue name")
	flag.StringVar(&config.MongoDB, "MongoDB", "localhost:27017", "MongoDB server addresses comma delimited")
	flag.StringVar(&config.Database, "Database", "syros", "MongoDB database name")
	flag.IntVar(&config.DatabaseStale, "DatabaseStale", 5, "Deletes database records older than specified value in minutes, set 0 to disable")
	flag.Parse()

	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	repo, err := NewRepository(config)
	if err != nil {
		log.Fatalf("MongoDB connection error %v", err)
	}

	repo.Initialize()
	log.Infof("Connected to MongoDB cluster %v database initialization done", config.MongoDB)

	repo.RunGarbageCollector([]string{"containers", "hosts", "checks", "syros_services", "vsphere_hosts", "vsphere_dstores", "vsphere_vms"})

	nc, err := NewNatsConnection(config.Nats)
	if err != nil {
		log.Fatalf("Nats connection error %v", err)
	}
	defer nc.Close()

	log.Infof("Connected to NATS server %v status %v", nc.Conn.ConnectedUrl(), nc.Conn.Status())

	registry := NewRegistry(config, nc, repo)
	registry.WatchForAgents()
	registry.Start()

	consumer, err := NewConsumer(config, nc, repo)
	if err != nil {
		log.Fatalf("Consumer init error %v", err)
	}
	consumer.Consume()

	server := &HttpServer{
		Config: config,
	}
	log.Infof("Starting HTTP server on port %v", config.Port)
	go server.Start()

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-sigChan
	log.Infof("Shutting down %v signal received", sig)
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
