package main

// Config holds global configuration, defaults are provided in main.
type Config struct {
	LogLevel       string `m:"LogLevel"`
	Port           int    `m:"Port"`
	Nats           string `m:"Nats"`
	CollectorTopic string `m:"CollectorTopic"`
	CollectorQueue string `m:"CollectorQueue"`
	RegistryTopic  string `m:"RegistryTopic"`
	RegistryQueue  string `m:"RegistryQueue"`
	RethinkDB      string `m:"RethinkDB"`
	Database       string `m:"Database"`
	DatabaseStale  int    `m:"DatabaseStale"`
}
