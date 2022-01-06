package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Api           Api
	Database      Database
	Cache         Cache
	Elasticsearch Elasticsearch
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Api:           API(),
		Database:      DataStore(),
		Cache:         NewCache(),
		Elasticsearch: ElasticSearch(),
	}
}
