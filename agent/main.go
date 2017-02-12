package main

import (
	"context"
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var config = &Config{}
	flag.StringVar(&config.Environment, "Environment", "DEBUG", "environment: DEBUG, DEV, TEST, STG, PROD")
	flag.StringVar(&config.LogLevel, "LogLevel", "debug", "logging threshold level: debug|info|warn|error|fatal|panic")
	flag.IntVar(&config.Port, "Port", 8000, "HTTP port to listen on")
	flag.Parse()

	setLogLevel(config.LogLevel)
	log.Infof("Config: %+v", config)

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient(os.Getenv("DOCKER_HOST"), "", nil, defaultHeaders)
	if err != nil {
		log.Fatal(err)
	}

	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		log.Fatal(err)
	}

	host, err := cli.Info(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Host %s ContainersRunning %v\n", host.Name, host.ContainersRunning)

	for _, container := range containers {
		containerInfo, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s %s\n", containerInfo.Name, container.State, container.Status)
	}

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	log.Info("Shutdown")
}

func setLogLevel(levelName string) {
	level, err := log.ParseLevel(levelName)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
