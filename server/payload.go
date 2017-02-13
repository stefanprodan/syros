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
}
