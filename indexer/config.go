package main

// Config holds global configuration, defaults are provided in main.
type Config struct {
	LogLevel       string `m:"LogLevel"`
	Port           int    `m:"Port"`
	Nats           string `m:"Nats"`
	CollectorQueue string `m:"CollectorQueue"`
	RegistryTopic  string `m:"RegistryTopic"`
	RegistryQueue  string `m:"RegistryQueue"`
	MongoDB        string `m:"MongoDB"`
	Database       string `m:"Database"`
	DatabaseStale  int    `m:"DatabaseStale"`
}
