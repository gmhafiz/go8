package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Elasticsearch struct {
	Address  string `default:"http://localhost:9200"`
	User     string
	Password string
}

func ElasticSearch() Elasticsearch {
	var elasticsearch Elasticsearch
	envconfig.MustProcess("ELASTICSEARCH", &elasticsearch)

	return elasticsearch
}
