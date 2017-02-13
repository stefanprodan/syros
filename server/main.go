package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/nats-io/go-nats"
	"os"
	"os/signal"
	"syscall"
)

type DockerPayload struct {
	Host              types.Info
	ContainerList     []types.Container
	ContainerInfoList []types.ContainerJSON
}

func main() {
	nc, err := nats.Connect("nats://localhost:4222",
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Warnf("Got disconnected from NATS %v", nc.ConnectedUrl())
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("Got reconnected to NATS %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Errorf("NATS connection closed. Reason: %q", nc.LastError())
		}),
	)
	if err != nil {
		log.Fatalf("Nats connection error %v", err)
	}
	log.Infof("Connected to NATS server %v status %v", nc.ConnectedUrl(), nc.Status())
	defer nc.Close()

	nc.QueueSubscribe("docker", "syros", func(m *nats.Msg) {
		var payload DockerPayload
		err := json.Unmarshal(m.Data, &payload)
		if err != nil {
			log.Errorf("DockerPayload unmarshal error %v", err)
		} else {
			log.Infof("Host %v running containes %v", payload.Host.Name, payload.Host.ContainersRunning)
		}
	})

	//wait for SIGINT (Ctrl+C) or SIGTERM (docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Infof("Shuting down %v signal received", sig)
}
