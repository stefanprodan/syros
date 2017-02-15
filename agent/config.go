package main

// Config holds global configuration, defaults are provided in main.
// Agent config is populated from startup flag.
type Config struct {
	LogLevel           string
	Port               int
	CollectInterval    int
	DockerApiAddresses string
	Nats               string
	CollectorTopic     string
}
