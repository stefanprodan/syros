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

	pgmon, err := NewPGMonitor(config.PostgresURI, status)
	if err != nil {
		log.Fatalf("Postgres connection failed %s", err.Error())
	}

	isMaster, err := pgmon.IsMaster()
	if err != nil {
		log.Fatalf("Can't determine Postgres cluster state %s", err.Error())
	} else {
		log.Infof("Postgres Master %v", isMaster)
	}
	go pgmon.Start()

	election, err := NewElection(config.ConsulURI, config.ConsulTTL, config.ConsulKV, config.Hostname, status)
	if err != nil {
		log.Fatal(err.Error())
	}

	leader, err := election.GetLeaderWithRetry(5, 1)
	if err != nil {
		//TODO: stop pg service if role is master
		log.Fatalf("Consul connection failed %s", err.Error())
	}

	if len(leader) > 0 {
		if isMaster {
			//TODO: stop pg service, this should never happen
			log.Warnf("Conflict detected: leader is %v but this pg node %v is master", leader, config.Hostname)
		} else {
			log.Infof("Leader is %v joining cluster as follower", leader)
		}
	} else {
		if isMaster {
			log.Infof("No leader found and this pg node %v is master, joining cluster as leader", config.Hostname)
		} else {
			//TODO: wait till the salve is up to date (do not enter election mode)
			log.Warnf("Conflict detected: no leader found but this pg node %v is slave", config.Hostname)
		}
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
	pgmon.Stop()
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
