package configs

import (
	"log"

	"github.com/joho/godotenv"
)

var (
	DockerPort = "5433"
)

type Configs struct {
	Api           *Api
	Database      *Database
	Cache         *Cache
	Elasticsearch *Elasticsearch
}

func New() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	api := API()
	dataStore := DataStore()
	cache := NewCache()
	es := ElasticSearch()

	return &Configs{
		Api:           api,
		Database:      dataStore,
		Cache:         cache,
		Elasticsearch: es,
	}
}
