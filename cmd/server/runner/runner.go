package runner

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/undercode99/article_service/internal/api"
	"github.com/undercode99/article_service/internal/app/article"
	"github.com/undercode99/article_service/internal/database"
	"github.com/undercode99/article_service/internal/searching"
	"gorm.io/gorm"
)

type AppRunner struct {
	db            *gorm.DB
	apiService    *api.ApiService
	elasticClient *elasticsearch.TypedClient
}

// NewAppRunner initializes and returns an instance of the AppRunner struct.
//
// Parameters:
// - db: a pointer to a gorm.DB object, the database connection.
// - apiService: a pointer to an api.ApiService object, the API service.
// - elasticClient: a pointer to an elasticsearch.TypedClient object, the Elasticsearch client.
//
// Returns:
// - a pointer to an AppRunner object.
func NewAppRunner(db *gorm.DB, apiService *api.ApiService, elasticClient *elasticsearch.TypedClient) *AppRunner {
	return &AppRunner{
		db:            db,
		apiService:    apiService,
		elasticClient: elasticClient,
	}
}

// Migrate migrates the database and creates an index in Elasticsearch.
//
// ctx - The context of the function.
// Returns an error if the migration or index creation fails.
func (a *AppRunner) Migrate(ctx context.Context) {
	log.Println("Migrating...")
	err := database.MigrateDatabase(a.db)
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// create index
	err = searching.CreateIndexElastic(ctx, a.elasticClient, article.IndexName)
	if err != nil {
		log.Fatalf("failed to create index: %v", err)
	}
}

// Run runs the AppRunner.
//
// It migrates the context and runs the apiService.
func (a *AppRunner) Run(ctx context.Context) {
	a.Migrate(ctx)
	a.apiService.Run(ctx)
}
