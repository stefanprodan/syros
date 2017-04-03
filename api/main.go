package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/goware/jwtauth"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	var config = &Config{}
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8888, "HTTP port to listen on")
	flag.StringVar(&config.MongoDB, "MongoDB", "localhost:27017", "MongoDB server addresses comma delimited")
	flag.StringVar(&config.Database, "Database", "syros", "MongoDB database name")
	flag.StringVar(&config.JwtSecret, "JwtSecret", "syros", "JWT secret")
	flag.StringVar(&config.Credentials, "Credentials", "admin@admin", "Credentials format user@password")
	flag.StringVar(&config.AppPath, "AppPath", "", "Path to dist dir")
	flag.Parse()

	setLogLevel(config.LogLevel)

	log.Infof("Starting with config: %+v", config)

	if config.AppPath == "" {
		workDir, _ := os.Getwd()
		config.AppPath = filepath.Join(workDir, "dist")
		if config.LogLevel != "debug" {
			if _, err := os.Stat(filepath.Join(config.AppPath, "index.html")); err != nil {
				if os.IsNotExist(err) {
					log.Fatalf("index.html not found in %v", config.AppPath)
				} else {
					log.Fatalf("Path to dist dir %v error %v", config.AppPath, err.Error())
				}
			}
		}
	}

	repo, err := NewRepository(config)
	if err != nil {
		log.Fatalf("MongoDB connection error %v", err)
	}

	server := HttpServer{
		Config:     config,
		Repository: repo,
		TokenAuth:  jwtauth.New("HS256", []byte(config.JwtSecret), nil),
	}

	log.Infof("Starting HTTP server on port %v", config.Port)
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
