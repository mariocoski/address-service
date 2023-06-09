package config

import (
	"os"
)

type Config struct {
	Environment           string
	PostgresConnectionUrl string
	RepositoryType        string
	SentryUrl             string
	Auth0Domain           string
	Auth0Audience         string
	LogLevel              string
	ApiPort               string
	MongoDBConnectionUrl  string
	MongoDBName           string
}

func NewConfig() *Config {
	return &Config{
		Environment:           os.Getenv("ENVIRONMENT"),
		PostgresConnectionUrl: os.Getenv("POSTGRES_CONNECTION_URL"),
		MongoDBConnectionUrl:  os.Getenv("MONGODB_CONNECTION_URL"),
		MongoDBName:           os.Getenv("MONGODB_NAME"),
		RepositoryType:        os.Getenv("REPOSITORY_TYPE"),
		SentryUrl:             os.Getenv("SENTRY_URL"),
		Auth0Domain:           os.Getenv("AUTH0_DOMAIN"),
		Auth0Audience:         os.Getenv("AUTH0_AUDIENCE"),
		LogLevel:              os.Getenv("LOG_LEVEL"),
		ApiPort:               os.Getenv("API_PORT"),
	}
}
