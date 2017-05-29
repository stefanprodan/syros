package main

// Config holds global configuration, defaults are provided in main.
type Config struct {
	Environment string `m:"Environment"`
	LogLevel    string `m:"LogLevel"`
	Port        int    `m:"Port"`
	Hostname    string `m:"Hostname"`
	ConsulURI   string `m:"ConsulURI"`
	ConsulTTL   string `m:"ConsulTTL"`
	ConsulKV    string `m:"ConsulKV"`
	PostgresURI string `m:"PostgresURI"`
	NatsURI     string `m:"NatsURI"`
}
