package main

// Config holds global configuration, defaults are provided in main.
type Config struct {
	LogLevel       string
	Port           int
	Nats           string
	CollectorTopic string
	CollectorQueue string
	RegistryTopic  string
	RegistryQueue  string
	RethinkDB      string
	Database       string
	DatabaseStale  int
}
