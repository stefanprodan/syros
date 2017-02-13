package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/nats-io/go-nats"
	"time"
)

type DockerCollector struct {
	Host     string
	Client   *docker.Client
	StopChan chan struct{}
}

type DockerPayload struct {
	Host              types.Info
	ContainerList     []types.Container
	ContainerInfoList []types.ContainerJSON
}

func NewNatsConnection(servers string) (*nats.Conn, error) {
	nc, err := nats.Connect(servers,
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Warnf("Got disconnected from NATS %v", servers)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Infof("Got reconnected to NATS %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Errorf("NATS connection closed. Reason: %q", nc.LastError())
		}),
	)
	return nc, err
}

func NewDockerCollector(host string) (*DockerCollector, error) {

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	client, err := docker.NewClient(host, "", nil, defaultHeaders)
	if err != nil {
		return nil, err
	}
	collector := &DockerCollector{
		Client:   client,
		Host:     host,
		StopChan: make(chan struct{}, 1),
	}

	return collector, nil
}

func (col *DockerCollector) Collect() (*DockerPayload, error) {
	start := time.Now().UTC()
	payload := &DockerPayload{}

	host, err := col.Client.Info(context.Background())
	if err != nil {
		return nil, err
	}
	payload.Host = host

	//fmt.Printf("Host %s ContainersRunning %v\n", host.Name, host.ContainersRunning)

	options := types.ContainerListOptions{All: true}
	containers, err := col.Client.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}
	payload.ContainerList = containers
	details := make([]types.ContainerJSON, len(containers))
	for _, container := range containers {
		containerInfo, err := col.Client.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Error(err)
			continue
		}
		details = append(details, containerInfo)

		//fmt.Printf("%s %s %s\n", containerInfo.Name, container.State, container.Status)
	}
	payload.ContainerInfoList = details

	log.Debugf("%v collect duration: %v", col.Host, time.Now().UTC().Sub(start))
	return payload, nil
}
