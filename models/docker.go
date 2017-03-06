package models

import "time"

type DockerPayload struct {
	Host       DockerHost        `json:"host"`
	Containers []DockerContainer `json:"containers"`
}

type EnvironmentStats struct {
	Environment       string `json:"environment"`
	Hosts             int    `json:"hosts"`
	ContainersRunning int    `json:"containers_running"`
	NCPU              int    `json:"ncpu"`
	MemTotal          int64  `json:"mem_total"`
}

type EnvironmentDto struct {
	Host        DockerHost        `json:"host"`
	Containers  []DockerContainer `json:"containers"`
	Deployments ChartDto          `json:"deployments"`
}

type ChartDto struct {
	Labels []string `json:"labels"`
	Values []int64  `json:"values"`
}

type DockerHost struct {
	Id                 string    `gorethink:"id,omitempty" json:"id"`
	Containers         int       `gorethink:"containers" json:"containers"`
	ContainersRunning  int       `gorethink:"containers_running" json:"containers_running"`
	ContainersPaused   int       `gorethink:"containers_paused" json:"containers_paused"`
	ContainersStopped  int       `gorethink:"containers_stopped" json:"containers_stopped"`
	Images             int       `gorethink:"images" json:"images"`
	Driver             string    `gorethink:"driver" json:"driver"`
	SystemTime         string    `gorethink:"system_time" json:"system_time"`
	LoggingDriver      string    `gorethink:"logging_driver" json:"logging_driver"`
	CgroupDriver       string    `gorethink:"cgroup_driver" json:"cgroup_driver"`
	NEventsListener    int       `gorethink:"n_events_listener" json:"n_events_listener"`
	KernelVersion      string    `gorethink:"kernel_version" json:"kernel_version"`
	OperatingSystem    string    `gorethink:"operating_system" json:"operating_system"`
	OSType             string    `gorethink:"os_type" json:"os_type"`
	Architecture       string    `gorethink:"architecture" json:"architecture"`
	IndexServerAddress string    `gorethink:"index_server_address" json:"index_server_address"`
	NCPU               int       `gorethink:"ncpu" json:"ncpu"`
	MemTotal           int64     `gorethink:"mem_total" json:"mem_total"`
	DockerRootDir      string    `gorethink:"docker_root_dir" json:"docker_root_dir"`
	HTTPProxy          string    `gorethink:"http_proxy" json:"http_proxy"`
	HTTPSProxy         string    `gorethink:"https_proxy" json:"https_proxy"`
	NoProxy            string    `gorethink:"no_proxy" json:"no_proxy"`
	Name               string    `gorethink:"name" json:"name"`
	Labels             []string  `gorethink:"labels" json:"labels"`
	ExperimentalBuild  bool      `gorethink:"experimental_build" json:"experimental_build"`
	ServerVersion      string    `gorethink:"server_version" json:"server_version"`
	ClusterStore       string    `gorethink:"cluster_store" json:"cluster_store"`
	ClusterAdvertise   string    `gorethink:"cluster_advertise" json:"cluster_advertise"`
	DefaultRuntime     string    `gorethink:"default_runtime" json:"default_runtime"`
	LiveRestoreEnabled bool      `gorethink:"live_restore_enabled" json:"live_restore_enabled"`
	Registries         []string  `gorethink:"registries" json:"registries"`
	Collected          time.Time `gorethink:"collected" json:"collected"`
	Environment        string    `gorethink:"environment" json:"environment"`
}

type DockerContainer struct {
	Id            string            `gorethink:"id,omitempty" json:"id"`
	HostId        string            `gorethink:"host_id" json:"host_id"`
	HostName      string            `gorethink:"host_name" json:"host_name"`
	Port          string            `gorethink:"port" json:"port"`
	Image         string            `gorethink:"image" json:"image"`
	Command       string            `gorethink:"command" json:"command"`
	Labels        map[string]string `gorethink:"labels" json:"labels"`
	State         string            `gorethink:"state" json:"state"`
	Status        string            `gorethink:"status" json:"status"`
	Created       time.Time         `gorethink:"created" json:"created"`
	Path          string            `gorethink:"path" json:"path"`
	Args          []string          `gorethink:"args" json:"args"`
	Name          string            `gorethink:"name" json:"name"`
	RestartCount  int               `gorethink:"restart_count" json:"restart_count"`
	Env           []string          `gorethink:"env" json:"env"`
	PortBindings  map[string]string `gorethink:"port_bindings" json:"port_bindings"`
	NetworkMode   string            `gorethink:"network_mode" json:"network_mode"`
	RestartPolicy string            `gorethink:"restart_policy" json:"restart_policy"`
	StartedAt     time.Time         `gorethink:"started_at" json:"started_at"`
	FinishedAt    time.Time         `gorethink:"finished_at" json:"finished_at"`
	ExitCode      int               `gorethink:"exit_code" json:"exit_code"`
	Error         string            `gorethink:"error" json:"error"`
	Collected     time.Time         `gorethink:"collected" json:"collected"`
	Environment   string            `gorethink:"environment" json:"environment"`
}
