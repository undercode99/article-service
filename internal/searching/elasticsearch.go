package searching

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/undercode99/article_service/config"
)

// NewElasticClient creates a new Elasticsearch client.
//
// It takes a context and a config as parameters.
// Returns a pointer to the elasticsearch.TypedClient.
func NewElasticClient(ctx context.Context, cfg *config.Config) *elasticsearch.TypedClient {

	// create a new client
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{cfg.ElasticUrl},
	})

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := client.Ping().Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if !res {
		log.Fatal("Elasticsearch is not reachable")
	}

	return client

}

// CreateIndexElastic creates an index in Elasticsearch.
//
// ctx: the context to use for the request.
// client: the Elasticsearch client.
// index: the name of the index to create.
// Returns an error if there was a problem creating the index.
func CreateIndexElastic(ctx context.Context, client *elasticsearch.TypedClient, index string) error {
	indexExists, err := client.Indices.Exists(index).Do(ctx)
	if err != nil {
		log.Printf("Error checking if index exists: %v", err)
		return err
	}

	if indexExists {
		log.Printf("Index %s already exists", index)
		return nil
	}

	_, err = client.Indices.Create(index).Do(ctx)
	if err != nil {
		log.Printf("Error creating index: %v", err)
		return err
	}

	log.Printf("Index %s created successfully", index)
	return nil
}
