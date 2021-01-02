package configs

import (
	"log"

	"github.com/joho/godotenv"
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
