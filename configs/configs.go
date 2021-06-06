package configs

import (
	"log"

	"github.com/joho/godotenv"
)

type Configs struct {
	Api           Api
	Database      Database
	Cache         Cache
	Elasticsearch Elasticsearch
}

func New() *Configs {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Configs{
		Api:           API(),
		Database:      DataStore(),
		Cache:         NewCache(),
		Elasticsearch: ElasticSearch(),
	}
}
