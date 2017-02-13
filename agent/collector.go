package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
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
