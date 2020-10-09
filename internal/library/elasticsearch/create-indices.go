package elasticsearch

import (
	"context"

	"github.com/olivere/elastic/v7"
)

// CreateIndices Create Indices by name and its JSON mapping
func CreateIndices(es *elastic.Client) error {
	// Mapping types: https://www.elastic./co/guide/en/elasticsearch/reference/current/mapping-types.html
	mapping := `{
				"settings":{
					"number_of_shards":1,
					"number_of_replicas":0
				},
				"mappings":{
					"properties":{
						"title":{
							"type":"text"
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
	err := createIndex(es, "go8-books", mapping)
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
							"type":"text"
						}
					}
				}
			}`
	err = createIndex(es, "go8-authors", mapping)
	if err != nil {
		return err
	}

	return nil
}

func createIndex(es *elastic.Client, indexName, mapping string) error {
	exists, err := es.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}

	ctx := context.Background()
	if !exists {

		createIndex, err := es.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return err
		}

		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return nil
}
