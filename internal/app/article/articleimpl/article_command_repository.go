package articleimpl

import (
	"context"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/undercode99/article_service/internal/app/article"
	"gorm.io/gorm"
)

type ArticleCommandRepository struct {
	db            *gorm.DB
	elasticClient *elasticsearch.TypedClient
}

func NewArticleCommandRepository(db *gorm.DB, elasticClient *elasticsearch.TypedClient) article.ArticleCommandRepository {
	return &ArticleCommandRepository{
		db:            db,
		elasticClient: elasticClient,
	}
}

// CreateArticle creates a new article in the ArticleCommandRepository.
//
// It takes an article object as a parameter and returns an error.
func (r *ArticleCommandRepository) CreateArticle(article *article.Article) error {
	// Create a new article record in the database.
	created := r.db.Create(article)

	// Check if there was an error while creating the article.
	if created.Error != nil {
		return created.Error
	}

	// Return nil (no error) if the article was created successfully.
	return nil
}

// CreateIndexArticle indexes the document and creates an index.
//
// ctx: the context.Context object for handling deadlines, cancellations, and values across API boundaries.
// item: a pointer to the article.Article object to be indexed.
// Returns an error if there was a problem indexing the document.
func (r *ArticleCommandRepository) CreateIndexArticle(ctx context.Context, item *article.Article) error {
	// Convert the item ID to a string
	idString := strconv.Itoa(item.ID)

	// Index the document using the ElasticSearch client
	_, err := r.elasticClient.Index(article.IndexName).Id(idString).Request(item).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
