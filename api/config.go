package main

// Config holds global configuration, defaults are provided in main.
type Config struct {
	LogLevel  string
	Port      int
	RethinkDB string
	Database  string
	JwtSecret string
}
