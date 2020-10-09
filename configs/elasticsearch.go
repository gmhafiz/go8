package configs

import (
	"os"
)


// Elasticsearch Holds information necessary to create a client
type Elasticsearch struct {
	Address  string
	Username string
	Password string
}

func NewElasticSearch() *Elasticsearch {
	return &Elasticsearch{
		Address:  os.Getenv("ELASTICSEARCH_ADDRESS"),
		Username: os.Getenv("ELASTICSEARCH_USER"),
		Password: os.Getenv("ELASTICSEARCH_PASS"),
	}
}