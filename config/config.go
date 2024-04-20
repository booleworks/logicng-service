package config

import "time"

type Config struct {
	Host                  string
	Port                  string
	SyncComputationTimout time.Duration
}

func Default() *Config {
	timeout, _ := time.ParseDuration("5s")
	return &Config{
		Host:                  "localhost",
		Port:                  "8080",
		SyncComputationTimout: timeout,
	}
}
