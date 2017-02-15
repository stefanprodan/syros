package models

type DockerPayload struct {
	Host       DockerHost
	Containers []DockerContainer
}

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
