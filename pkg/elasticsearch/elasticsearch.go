package elasticsearch

import (
	"context"

	"github.com/olivere/elastic/v7"
)

// Config Holds information necessary to create a client
type Config struct {
	Address  string `yaml:"ADDRESS"`
	Username string `yaml:"USERNAME"`
	Password string `yaml:"PASSWORD"`
}

// Es Holds The ElasticSearch client and config
type Es struct {
	Client *elastic.Client
	Config *Config
}

// New Create a new ElasticSearch instance
func New(esConfig Config) (*Es, error) {
	client, err := elastic.NewClient(elastic.SetURL(esConfig.Address))
	if err != nil {
		return nil, err
	}

	return &Es{
		Client: client,
	}, nil
}

// CreateIndices Create Indices by name and its JSON mapping
func (es *Es) CreateIndices() error {
	// Mapping types: https://www.elastic./co/guide/en/elasticsearch/reference/current/mapping-types.html
	mapping := `{
				"settings":{
					"number_of_shards":1,
					"number_of_replicas":0
				},
				"mappings":{
					"properties":{
						"title":{
							"type":"completion"
						},
						"published_date":{
							"type":"text"
						},
						"description":{
							"type":"text"
						},
						"image_url":{
							"type":"text"
						}
					}
				}
			}`
	err := es.createIndex("go8-books", mapping)
	if err != nil {
		return err
	}

	// Mapping types: https://www.elastic./co/guide/en/elasticsearch/reference/current/mapping-types.html
	mapping = `{
				"settings":{
					"number_of_shards":1,
					"number_of_replicas":0
				},
				"mappings":{
					"properties":{
						"first_name":{
							"type":"text"
						},
						"middle_name":{
							"type":"text"
						},
						"last_name":{
							"type":"completion"
						}
					}
				}
			}`
	err = es.createIndex("go8-authors", mapping)
	if err != nil {
		return err
	}

	return nil
}

func (es *Es) createIndex(indexName, mapping string) error {
	exists, err := es.Client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}

	ctx := context.Background()
	if !exists {

		createIndex, err := es.Client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return nil
}
