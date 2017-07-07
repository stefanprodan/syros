package main

// config holds global configuration, defaults are provided in main.
type Config struct {
	Environment   string `m:"Environment"`
	LogLevel      string `m:"LogLevel"`
	Port          int    `m:"Port"`
	Hostname      string `m:"Hostname"`
	ConsulURI     string `m:"ConsulURI"`
	ConsulTTL     string `m:"ConsulTTL"`
	ConsulKV      string `m:"ConsulKV"`
	ConsulRetry   int    `m:"ConsulRetry"`
	PostgresURI   string `m:"PostgresURI"`
	PostgresCheck int    `m:"PostgresCheck"`
	NatsURI       string `m:"NatsURI"`
	User          string `m:"User"`
}
