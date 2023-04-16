package config

import "os"

type Config struct {
	Environment           string
	PostgresConnectionUrl string
	RepositoryType        string
	SentryUrl             string
}

func NewConfig() *Config {
	return &Config{
		Environment:           os.Getenv("ENVIRONMENT"),
		PostgresConnectionUrl: os.Getenv("POSTGRES_CONNECTION_URI"),
		RepositoryType:        os.Getenv("REPOSITORY_TYPE"),
		SentryUrl:             os.Getenv("SENTRY_URL"),
	}
}
