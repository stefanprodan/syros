package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var config = &Config{}
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8000, "HTTP port to listen on")
	flag.StringVar(&config.Nats, "Nats", "nats://localhost:4222", "Nats server addresses comma delimited")
	flag.StringVar(&config.CollectorTopic, "CollectorTopic", "docker", "Nats collector topic name")
	flag.StringVar(&config.CollectorQueue, "CollectorQueue", "syros", "Nats collector queue name")
	flag.StringVar(&config.RethinkDB, "RethinkDB", "localhost:28015", "RethinkDB server addresses comma delimited")
	flag.StringVar(&config.Database, "Database", "syros", "RethinkDB database name")
	flag.Parse()

	repo, err := NewRepository(config)
	if err != nil {
		log.Fatalf("RethinkDB connection error %v", err)
	}

	repo.Initialize()
	log.Infof("Connected to RethinkDB cluster %v database initialization done", config.RethinkDB)

	nc, err := NewNatsConnection(config.Nats)
	if err != nil {
		log.Fatalf("Nats connection error %v", err)
	}
	defer nc.Close()

	log.Infof("Connected to NATS server %v status %v", nc.ConnectedUrl(), nc.Status())

	consumer, err := NewDockerConsumer(config, nc, repo)
	if err != nil {
		log.Fatalf("Docker consumer init error %v", err)
	}
	consumer.Consume()

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Infof("Shuting down %v signal received", sig)
}
