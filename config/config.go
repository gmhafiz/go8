package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	API
	Cors

	Database
	Cache
	Elasticsearch

	OpenTelemetry
	Session
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	return &Config{
		API:           NewAPI(),
		Cors:          NewCors(),
		Database:      DataStore(),
		Cache:         NewCache(),
		Elasticsearch: ElasticSearch(),
		Session:       NewSession(),
		OpenTelemetry: NewOpenTelemetry(),
	}
}
