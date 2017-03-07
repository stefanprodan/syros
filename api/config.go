package main

// Config holds global configuration, defaults are provided in main.
type Config struct {
	LogLevel    string `m:"LogLevel"`
	Port        int    `m:"Port"`
	RethinkDB   string `m:"RethinkDB"`
	Database    string `m:"Database"`
	JwtSecret   string `m:"JwtSecret"`
	Credentials string `m:"Credentials"`
	AppPath     string `m:"AppPath"`
}
