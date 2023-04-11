package config

import "os"

type Config struct {
	Environment           string
	PostgresConnectionUrl string
}

func NewConfig() *Config {
	return &Config{
		Environment:           os.Getenv("ENVIRONMENT"),
		PostgresConnectionUrl: os.Getenv("POSTGRES_CONNECTION_URI"),
	}
}
