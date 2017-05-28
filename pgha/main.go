package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

func main() {

	election, err := NewElection("localhost:8500", "pgha/leader/election", "pgdev")
	if err != nil {
		panic(err)
	}

	go election.Start()

	//wait for exit signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-sigChan
	log.Infof("Shutting down %v signal received", sig)
	election.Stop()
}
