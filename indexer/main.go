package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
	flag.IntVar(&config.DatabaseStaleSince, "DatabaseStaleSince", 48, "Scan for database records since specified value in hours")
	flag.Parse()

	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	repo, err := NewMongoRepository(config)
	if err != nil {
		log.Fatalf("MongoDB connection error %v", err)
	}

	repo.Initialize()
	log.Infof("Connected to MongoDB cluster %v database initialization done", config.MongoDB)

	repo.RunGarbageCollector([]string{"containers", "hosts", "checks", "syros_services"})

	nc, err := NewNatsConnection(config.Nats)
	if err != nil {
		log.Fatalf("Nats connection error %v", err)
	}
	defer nc.Close()

	log.Infof("Connected to NATS server %v status %v", nc.ConnectedUrl(), nc.Status())

	registry := NewRegistry(config, nc, repo)
	registry.WatchForAgents()

	indexer := models.SyrosService{
		Environment: "all",
		Type:        "indexer",
	}
	indexer.Config, _ = models.ConfigToMap(config, "m")
	indexer.Hostname, _ = os.Hostname()
	uuid, _ := models.NewUUID()
	indexer.Id = models.Hash(indexer.Hostname + uuid)
	go func(a models.SyrosService) {
		for true {
			indexer.Collected = time.Now().UTC()
			registry.SelfRegister(indexer)
			time.Sleep(10 * time.Second)
		}
	}(indexer)

	consumer, err := NewConsumer(config, nc, repo)
	if err != nil {
		log.Fatalf("Consumer init error %v", err)
	}
	consumer.Consume()

	//server := &HttpServer{
	//	Config:     config,
	//	Registry:   registry,
	//	Repository: repo,
	//}
	//log.Infof("Starting HTTP server on port %v", config.Port)
	//go server.Start()

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Infof("Shuting down %v signal received", sig)
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
