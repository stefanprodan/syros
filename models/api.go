package models

type EnvironmentStats struct {
	Environment       string `json:"environment" bson:"_id"`
	Hosts             int    `json:"hosts" bson:"hosts"`
	ContainersRunning int    `json:"containers_running" bson:"containers_running"`
	NCPU              int    `json:"ncpu" bson:"ncpu"`
	MemTotal          int64  `json:"mem_total" bson:"mem_total"`
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
