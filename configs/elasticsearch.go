package configs

import "os"

type Elasticsearch struct {
	Address  string
	User     string
	Password string
}

func ElasticSearch() *Elasticsearch {
	return &Elasticsearch{
		Address:  os.Getenv("ELASTICSEARCH_ADDRESS"),
		User:     os.Getenv("ELASTICSEARCH_USER"),
		Password: os.Getenv("ELASTICSEARCH_PASS"),
	}
}
