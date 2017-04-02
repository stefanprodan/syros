package models

import "time"

type DockerPayload struct {
	Host       DockerHost        `json:"host"`
	Containers []DockerContainer `json:"containers"`
}

type DockerHost struct {
	Id                 string    `bson:"_id,omitempty" json:"id"`
	Containers         int       `bson:"containers" json:"containers"`
	ContainersRunning  int       `bson:"containers_running" json:"containers_running"`
	ContainersPaused   int       `bson:"containers_paused" json:"containers_paused"`
	ContainersStopped  int       `bson:"containers_stopped" json:"containers_stopped"`
	Images             int       `bson:"images" json:"images"`
	Driver             string    `bson:"driver" json:"driver"`
	SystemTime         string    `bson:"system_time" json:"system_time"`
	LoggingDriver      string    `bson:"logging_driver" json:"logging_driver"`
	CgroupDriver       string    `bson:"cgroup_driver" json:"cgroup_driver"`
	NEventsListener    int       `bson:"n_events_listener" json:"n_events_listener"`
	KernelVersion      string    `bson:"kernel_version" json:"kernel_version"`
	OperatingSystem    string    `bson:"operating_system" json:"operating_system"`
	OSType             string    `bson:"os_type" json:"os_type"`
	Architecture       string    `bson:"architecture" json:"architecture"`
	IndexServerAddress string    `bson:"index_server_address" json:"index_server_address"`
	NCPU               int       `bson:"ncpu" json:"ncpu"`
	MemTotal           int64     `bson:"mem_total" json:"mem_total"`
	DockerRootDir      string    `bson:"docker_root_dir" json:"docker_root_dir"`
	HTTPProxy          string    `bson:"http_proxy" json:"http_proxy"`
	HTTPSProxy         string    `bson:"https_proxy" json:"https_proxy"`
	NoProxy            string    `bson:"no_proxy" json:"no_proxy"`
	Name               string    `bson:"name" json:"name"`
	Labels             []string  `bson:"labels" json:"labels"`
	ExperimentalBuild  bool      `bson:"experimental_build" json:"experimental_build"`
	ServerVersion      string    `bson:"server_version" json:"server_version"`
	ClusterStore       string    `bson:"cluster_store" json:"cluster_store"`
	ClusterAdvertise   string    `bson:"cluster_advertise" json:"cluster_advertise"`
	DefaultRuntime     string    `bson:"default_runtime" json:"default_runtime"`
	LiveRestoreEnabled bool      `bson:"live_restore_enabled" json:"live_restore_enabled"`
	Registries         []string  `bson:"registries" json:"registries"`
	Collected          time.Time `bson:"collected" json:"collected"`
	Environment        string    `bson:"environment" json:"environment"`
}

type DockerContainer struct {
	Id            string            `bson:"_id,omitempty" json:"id"`
	HostId        string            `bson:"host_id" json:"host_id"`
	HostName      string            `bson:"host_name" json:"host_name"`
	Port          string            `bson:"port" json:"port"`
	Image         string            `bson:"image" json:"image"`
	Command       string            `bson:"command" json:"command"`
	Labels        map[string]string `bson:"labels" json:"labels"`
	State         string            `bson:"state" json:"state"`
	Status        string            `bson:"status" json:"status"`
	Created       time.Time         `bson:"created" json:"created"`
	Path          string            `bson:"path" json:"path"`
	Args          []string          `bson:"args" json:"args"`
	Name          string            `bson:"name" json:"name"`
	RestartCount  int               `bson:"restart_count" json:"restart_count"`
	Env           []string          `bson:"env" json:"env"`
	PortBindings  map[string]string `bson:"port_bindings" json:"port_bindings"`
	NetworkMode   string            `bson:"network_mode" json:"network_mode"`
	RestartPolicy string            `bson:"restart_policy" json:"restart_policy"`
	StartedAt     time.Time         `bson:"started_at" json:"started_at"`
	FinishedAt    time.Time         `bson:"finished_at" json:"finished_at"`
	ExitCode      int               `bson:"exit_code" json:"exit_code"`
	Error         string            `bson:"error" json:"error"`
	Collected     time.Time         `bson:"collected" json:"collected"`
	Environment   string            `bson:"environment" json:"environment"`
}
