package configs

import (
	"github.com/kelseyhightower/envconfig"
)

type Elasticsearch struct {
	Address  string
	User     string
	Password string
}

func ElasticSearch() Elasticsearch {
	var elasticsearch Elasticsearch
	envconfig.MustProcess("ELASTICSEARCH", &elasticsearch)

	return elasticsearch
}
