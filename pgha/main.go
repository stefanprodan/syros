package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"time"

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
	flag.IntVar(&config.ConsulRetry, "ConsulRetry", 10, "Number of Consul connection reties")
	flag.StringVar(&config.PostgresURI, "PostgresURI", "postgres://user:password@localhost/db?sslmode=disable", "Postgres URI")
	flag.StringVar(&config.NatsURI, "NatsURI", "nats://localhost:4222", "Nats URI")
	flag.StringVar(&config.User, "User", "postgres", "User to run under")
	flag.Parse()
	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	if config.Hostname == "" {
		config.Hostname, _ = os.Hostname()
	}

	if config.Environment != "debug" {
		//check if running under postgres
		id := execId(10)
		if id != config.User {
			log.Fatalf("Running under %s expected postgres", id)
		}

		//check if repmgr is installed
		execRepmgrVersion(2)
	}

	status := NewStatus()

	pgmon, err := NewPGMonitor(config.PostgresURI, status)
	if err != nil {
		log.Fatalf("Postgres connection failed %s", err.Error())
	}

	isMaster, err := pgmon.IsMaster()
	if err != nil {
		log.Fatalf("Can't determine Postgres cluster state %s", err.Error())
	}
	status.SetPostgresStatus(isMaster)
	go pgmon.Start()

	pgstats, err := NewPGStats(config)
	if err != nil {
		log.Fatalf("PGStats init failed %s", err.Error())
	}

	stats, err := pgstats.GetReplicationStats()
	if err != nil {
		log.Fatalf("PGStats GetReplicationStats failed %s", err.Error())
	}
	log.Infof("Postgres replication stats: %+v", stats)

	err = pgstats.SaveReplicationStats(stats)
	if err != nil {
		log.Fatalf("PGStats SaveReplicationStats failed %s", err.Error())
	}

	election, err := NewElection(config, status)
	if err != nil {
		log.Fatal(err.Error())
	}

INIT:
	leader, err := election.GetLeaderWithRetry(10, 5)
	if err != nil {
		//stop pg service if is master
		if isMaster {
			execPgStop(20)
			log.Fatal("Stopping postgres service: this pg node is master but Consul is unreachable")
		}
		log.Fatalf("Consul connection failed %s", err.Error())
	}

	if len(leader) > 0 {
		if isMaster {
			//stop pg service, this should never happen
			execPgStop(20)
			log.Fatalf("Stopping postgres service: leader is %v but this pg node %v is master", leader, config.Hostname)
		} else {
			log.Infof("Should join cluster as follower since leader is %v", leader)
		}
	} else {
		if isMaster {
			log.Infof("Should join cluster as leader since no leader found and this pg node %v is master", config.Hostname)
		} else {
			//retry till there is a master running (do not enter election mode)
			log.Warnf("Conflict detected: no leader found but this pg node %v is slave, will wait 5sec and retry", config.Hostname)
			time.Sleep(5 * time.Second)
			//TODO: better use recursion
			goto INIT
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
