package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/stefanprodan/syros/models"
	"time"
)

type DockerCollector struct {
	ApiAddress string
	Client     *docker.Client
	Config     *Config
	StopChan   chan struct{}
}

func NewDockerCollector(apiAddress string, config *Config) (*DockerCollector, error) {

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	client, err := docker.NewClient(apiAddress, "", nil, defaultHeaders)
	if err != nil {
		return nil, err
	}
	collector := &DockerCollector{
		ApiAddress: apiAddress,
		Client:     client,
		Config:     config,
		StopChan:   make(chan struct{}, 1),
	}

	return collector, nil
}

func (col *DockerCollector) Collect() (*models.DockerPayload, error) {
	start := time.Now().UTC()
	payload := &models.DockerPayload{}

	host, err := col.Client.Info(context.Background())
	if err != nil {
		return nil, err
	}
	payload.Host = MapDockerHost(col.Config.Environment, host)

	options := types.ContainerListOptions{All: true}
	containers, err := col.Client.ContainerList(context.Background(), options)
	if err != nil {
		return nil, err
	}

	payload.Containers = make([]models.DockerContainer, 0)

	for _, container := range containers {
		containerInfo, err := col.Client.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			log.Error(err)
			continue
		}
		payload.Containers = append(payload.Containers, MapDockerContainer(col.Config.Environment, host.ID, host.Name, container, containerInfo))
	}

	log.Debugf("%v collect duration: %v containers %v", col.ApiAddress, time.Now().UTC().Sub(start), len(payload.Containers))
	return payload, nil
}

func MapDockerHost(environment string, info types.Info) models.DockerHost {
	host := models.DockerHost{
		Id:                 info.ID,
		Containers:         info.Containers,
		ContainersRunning:  info.ContainersRunning,
		ContainersPaused:   info.ContainersPaused,
		ContainersStopped:  info.ContainersStopped,
		Images:             info.Images,
		Driver:             info.Driver,
		SystemTime:         info.SystemTime,
		LoggingDriver:      info.LoggingDriver,
		CgroupDriver:       info.CgroupDriver,
		NEventsListener:    info.NEventsListener,
		KernelVersion:      info.KernelVersion,
		OperatingSystem:    info.OperatingSystem,
		OSType:             info.OSType,
		Architecture:       info.Architecture,
		IndexServerAddress: info.IndexServerAddress,
		NCPU:               info.NCPU,
		MemTotal:           info.MemTotal,
		DockerRootDir:      info.DockerRootDir,
		HTTPProxy:          info.HTTPProxy,
		HTTPSProxy:         info.HTTPSProxy,
		NoProxy:            info.NoProxy,
		Name:               info.Name,
		Labels:             info.Labels,
		ExperimentalBuild:  info.ExperimentalBuild,
		ServerVersion:      info.ServerVersion,
		ClusterStore:       info.ClusterStore,
		ClusterAdvertise:   info.ClusterAdvertise,
		DefaultRuntime:     info.DefaultRuntime,
		LiveRestoreEnabled: info.LiveRestoreEnabled,
		Collected:          time.Now().UTC(),
		Environment:        environment,
	}
	for _, reg := range info.RegistryConfig.IndexConfigs {
		host.Registries = append(host.Registries, reg.Name)
	}

	return host
}

func MapDockerContainer(environment string, hostId string, hostName string, c types.Container, cj types.ContainerJSON) models.DockerContainer {
	container := models.DockerContainer{
		Id:           c.ID,
		HostId:       hostId,
		HostName:     hostName,
		Image:        c.Image,
		Command:      c.Command,
		Labels:       c.Labels,
		State:        c.State,
		Status:       c.Status,
		Path:         cj.ContainerJSONBase.Path,
		Args:         cj.ContainerJSONBase.Args,
		Name:         cj.ContainerJSONBase.Name,
		RestartCount: cj.ContainerJSONBase.RestartCount,
		PortBindings: make(map[string]string),
		Collected:    time.Now().UTC(),
		Environment:  environment,
	}

	container.Created, _ = time.Parse(time.RFC3339, cj.ContainerJSONBase.Created)
	if len(container.Name) > 1 {
		container.Name = container.Name[1:len(container.Name)]
	}

	if cj.Config != nil {
		container.Env = cj.Config.Env
	}

	if cj.ContainerJSONBase.State != nil {

		container.StartedAt, _ = time.Parse(time.RFC3339, cj.ContainerJSONBase.State.StartedAt)
		container.FinishedAt, _ = time.Parse(time.RFC3339, cj.ContainerJSONBase.State.FinishedAt)
		container.ExitCode = cj.ContainerJSONBase.State.ExitCode
		container.Error = cj.ContainerJSONBase.State.Error
	}

	if cj.ContainerJSONBase.HostConfig != nil {
		container.NetworkMode = string(cj.ContainerJSONBase.HostConfig.NetworkMode)
		container.RestartPolicy = cj.ContainerJSONBase.HostConfig.RestartPolicy.Name
		for key, val := range cj.ContainerJSONBase.HostConfig.PortBindings {
			if len(val) > 0 {
				container.PortBindings[string(key)] = val[0].HostPort
			}
		}
	}

	// use first Host port bind as container Port
	if len(container.PortBindings) > 0 {
		for _, val := range container.PortBindings {
			if len(val) > 0 {
				container.Port = val
			}
		}
	}

	return container
}
