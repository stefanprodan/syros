package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/goware/jwtauth"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var config = &Config{}
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8888, "HTTP port to listen on")
	flag.StringVar(&config.RethinkDB, "RethinkDB", "localhost:28015", "RethinkDB server addresses comma delimited")
	flag.StringVar(&config.Database, "Database", "syros", "RethinkDB database name")
	flag.StringVar(&config.JwtSecret, "JwtSecret", "syros", "JWT secret")
	flag.StringVar(&config.Credentials, "Credentials", "admin@admin", "Credentials format user@password")
	flag.Parse()

	setLogLevel(config.LogLevel)
	log.Infof("Starting with config: %+v", config)

	repo, err := NewRepository(config)
	if err != nil {
		log.Fatalf("RethinkDB connection error %v", err)
	}

	server := HttpServer{
		Config:     config,
		Repository: repo,
		TokenAuth:  jwtauth.New("HS256", []byte(config.JwtSecret), nil),
	}

	go server.Start()

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
