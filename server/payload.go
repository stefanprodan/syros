package main

import (
	//log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
)

type DockerHost struct {
	Id                 string `gorethink:"id,omitempty"`
	Containers         int
	ContainersRunning  int
	ContainersPaused   int
	ContainersStopped  int
	Images             int
	Driver             string
	SystemTime         string
	LoggingDriver      string
	CgroupDriver       string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	OSType             string
	Architecture       string
	IndexServerAddress string
	NCPU               int
	MemTotal           int64
	DockerRootDir      string
	HTTPProxy          string
	HTTPSProxy         string
	NoProxy            string
	Name               string
	Labels             []string
	ExperimentalBuild  bool
	ServerVersion      string
	ClusterStore       string
	ClusterAdvertise   string
	DefaultRuntime     string
	LiveRestoreEnabled bool
	Registries         []string
}

func MapDockerHost(info types.Info) (DockerHost, error) {
	host := DockerHost{
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
	}
	for _, reg := range info.RegistryConfig.IndexConfigs {
		host.Registries = append(host.Registries, reg.Name)
	}

	return host, nil
}

type DockerContainer struct {
	Id            string `gorethink:"id,omitempty"`
	HostId        string `gorethink:"host_id,omitempty"`
	HostName      string
	Image         string // Container
	Command       string
	Labels        map[string]string
	State         string
	Status        string
	Created       string // ContainerJSON
	Path          string
	Args          []string
	Name          string
	RestartCount  int
	Env           []string          // ContainerJSON -> Config
	PortBindings  map[string]string // ContainerJSON -> HostConfig
	NetworkMode   string
	RestartPolicy string
	StartedAt     string // ContainerJSON -> State
	FinishedAt    string
	ExitCode      int
	Error         string
}

func MapDockerContiner(hostId string, c types.Container, cj types.ContainerJSON) (DockerContainer, error) {
	container := DockerContainer{
		Id:           c.ID,
		HostId:       hostId,
		Image:        c.Image,
		Command:      c.Command,
		Labels:       c.Labels,
		State:        c.State,
		Status:       c.Status,
		Created:      cj.ContainerJSONBase.Created,
		Path:         cj.ContainerJSONBase.Path,
		Args:         cj.ContainerJSONBase.Args,
		Name:         cj.ContainerJSONBase.Name,
		RestartCount: cj.ContainerJSONBase.RestartCount,
		PortBindings: make(map[string]string),
	}

	container.Name = container.Name[1:len(container.Name)]

	if cj.Config != nil {
		container.Env = cj.Config.Env
	}

	if cj.ContainerJSONBase.State != nil {
		container.StartedAt = cj.ContainerJSONBase.State.StartedAt
		container.FinishedAt = cj.ContainerJSONBase.State.FinishedAt
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

	return container, nil
}
