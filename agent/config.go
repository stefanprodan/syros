package main

// Config holds global configuration, defaults are provided in main.
// Agent config is populated from startup flag.
type Config struct {
	Environment            string `m:"Environment"`
	LogLevel               string `m:"LogLevel"`
	Port                   int    `m:"Port"`
	CollectInterval        int    `m:"CollectInterval"`
	DockerApiAddresses     string `m:"DockerApiAddresses"`
	ConsulApiAddresses     string `m:"ConsulApiAddresses"`
	VSphereApiAddress      string `m:"VSphereApiAddresses"`
	VSphereInclude         string `m:"VSphereInclude"`
	VSphereExclude         string `m:"VSphereExclude"`
	VSphereCollectInterval int    `m:"VSphereCollectInterval"`
	Nats                   string `m:"Nats"`
}
