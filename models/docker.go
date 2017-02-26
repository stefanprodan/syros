package models

import "time"

type DockerPayload struct {
	Host       DockerHost
	Containers []DockerContainer
}

type DockerHost struct {
	Id                 string    `gorethink:"id,omitempty" json:"id"`
	Containers         int       `json:"containers"`
	ContainersRunning  int       `json:"containers_running"`
	ContainersPaused   int       `json:"containers_paused"`
	ContainersStopped  int       `json:"containers_stopped"`
	Images             int       `json:"images"`
	Driver             string    `json:"driver"`
	SystemTime         string    `json:"system_time"`
	LoggingDriver      string    `json:"logging_driver"`
	CgroupDriver       string    `json:"cgroup_driver"`
	NEventsListener    int       `json:"n_events_listener"`
	KernelVersion      string    `json:"kernel_version"`
	OperatingSystem    string    `json:"operating_system"`
	OSType             string    `json:"os_type"`
	Architecture       string    `json:"architecture"`
	IndexServerAddress string    `json:"index_server_address"`
	NCPU               int       `json:"ncpu"`
	MemTotal           int64     `json:"mem_total"`
	DockerRootDir      string    `json:"docker_root_dir"`
	HTTPProxy          string    `json:"http_proxy"`
	HTTPSProxy         string    `json:"https_proxy"`
	NoProxy            string    `json:"no_proxy"`
	Name               string    `json:"name"`
	Labels             []string  `json:"labels"`
	ExperimentalBuild  bool      `json:"experimental_build"`
	ServerVersion      string    `json:"server_version"`
	ClusterStore       string    `json:"cluster_store"`
	ClusterAdvertise   string    `json:"cluster_advertise"`
	DefaultRuntime     string    `json:"default_runtime"`
	LiveRestoreEnabled bool      `json:"live_restore_enabled"`
	Registries         []string  `json:"registries"`
	Collected          time.Time `json:"collected"`
	Environment        string    `json:"environment"`
}

type DockerContainer struct {
	Id            string            `gorethink:"id,omitempty" json:"id"`
	HostId        string            `gorethink:"host_id,omitempty" json:"host_id"`
	HostName      string            `json:"host_name"`
	Port          string            `json:"port"`
	Image         string            `json:"image"` // Container
	Command       string            `json:"command"`
	Labels        map[string]string `json:"labels"`
	State         string            `json:"state"`
	Status        string            `json:"status"`
	Created       time.Time         `json:"created"` // ContainerJSON
	Path          string            `json:"path"`
	Args          []string          `json:"args"`
	Name          string            `json:"name"`
	RestartCount  int               `json:"restart_count"`
	Env           []string          `json:"env"`           // ContainerJSON -> Config
	PortBindings  map[string]string `json:"port_bindings"` // ContainerJSON -> HostConfig
	NetworkMode   string            `json:"network_mode"`
	RestartPolicy string            `json:"restart_policy"`
	StartedAt     time.Time         `json:"started_at"` // ContainerJSON -> State
	FinishedAt    time.Time         `json:"finished_at"`
	ExitCode      int               `json:"exit_code"`
	Error         string            `json:"error"`
	Collected     time.Time         `json:"collected"`
	Environment   string            `json:"environment"`
}
