package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

var version = "undefined"

func main() {
	var config = &Config{}
	flag.StringVar(&config.Environment, "Environment", "dev", "Environment dev|int|stg|test|prep|prod")
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 9898, "HTTP port to listen on")
	flag.StringVar(&config.Hostname, "Hostname", "", "Hostname")
	flag.StringVar(&config.ConsulURI, "ConsulURI", "localhost:8500", "Consul address")
	flag.StringVar(&config.ConsulTTL, "ConsulTTL", "10s", "Consul session TTL")
	flag.StringVar(&config.ConsulKV, "ConsulKV", "pgha", "Consul KV prefix")
	flag.StringVar(&config.PostgresURI, "PostgresURI", "postgres://user:password@localhost/db?sslmode=disable", "Postgres URI")
	flag.StringVar(&config.NatsURI, "NatsURI", "nats://localhost:4222", "Nats URI")
	flag.Parse()
	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	if config.Hostname == "" {
		config.Hostname, _ = os.Hostname()
	}

	status := NewStatus()

	election, err := NewElection(config.ConsulURI, config.ConsulTTL, config.ConsulKV, config.Hostname, status)
	if err != nil {
		log.Fatal(err.Error())
	}
	go election.Start()

	server, err := NewHttpServer(config, status, election)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Infof("Starting HTTP server on port %v", config.Port)
	go server.Start()

	//wait for exit signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-sigChan
	log.Infof("Shutting down %v signal received", sig)
	election.Stop()
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
