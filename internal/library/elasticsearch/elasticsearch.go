package elasticsearch

import (
	"log"

	"github.com/olivere/elastic/v7"

	"go8ddd/configs"
)

// Es Holds The ElasticSearch client and config
type Es struct {
	Client *elastic.Client
	Config *configs.Elasticsearch
}

func New(cfg *configs.Elasticsearch) *Es {
	client, err := elastic.NewClient(elastic.SetURL(cfg.Address))
	if err != nil {
		log.Panicf("elasticsearch error, %v", err.Error())
	}

	err = CreateIndices(client)
	if err != nil {
		log.Panicf("error creating indices, %v", err.Error())
	}

	return &Es{
		Client: client,
		Config: cfg,
	}
}

