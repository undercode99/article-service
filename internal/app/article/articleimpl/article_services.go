package articleimpl

import (
	"context"
	"log"

	"gorm.io/gorm"

	"github.com/undercode99/article_service/internal/app/article"
)

type ArticleService struct {
	articleCommandRepository article.ArticleCommandRepository
	articleQueryRepository   article.ArticleQueryRepository
	articleCachingRepository article.ArticleCachingRepository
}

// NewArticleService creates a new instance of the ArticleService struct.
//
// Parameters:
// - articleCommandRepository: an instance of the ArticleCommandRepository interface.
// - articleQueryRepository: an instance of the ArticleQueryRepository interface.
//
// Returns:
// - a pointer to the newly created ArticleService struct.
func NewArticleService(
	articleCommandRepository article.ArticleCommandRepository,
	articleQueryRepository article.ArticleQueryRepository,
	articleCachingRepository article.ArticleCachingRepository,
) article.ArticleService {
	return &ArticleService{
		articleCommandRepository: articleCommandRepository,
		articleQueryRepository:   articleQueryRepository,
		articleCachingRepository: articleCachingRepository,
	}
}

// CreateArticle creates a new article.
// It takes an article create command as a parameter and returns the created article and any error encountered.
func (s *ArticleService) CreateArticle(ctx context.Context, cmd *article.ArticleCreateCommand) (*article.Article, error) {
	// Validate the article create command
	if err := cmd.Validate(); err != nil {
		return nil, article.ErrArticleValidation
	}

	// Create the article using the article create command
	createdArticle := article.NewArticle(cmd)

	// Save the created article using the article command repository
	err := s.articleCommandRepository.CreateArticle(createdArticle)
	if err != nil {
		return nil, err
	}

	// Create the index for the article asynchronously
	go func() {
		err := s.articleCommandRepository.CreateIndexArticle(ctx, createdArticle)
		if err != nil {
			log.Printf("failed to create index for article: %v", err)
		}
	}()

	// Return the created article and no error
	return createdArticle, nil
}

// GetArticleByID retrieves an article by its ID.
// It takes a context.Context and an integer ID as parameters.
// It returns a pointer to an article.Article struct and an error.
func (s *ArticleService) GetArticleByID(ctx context.Context, id int) (*article.Article, error) {
	// Get the article from the cache
	articleCache, err := s.articleCachingRepository.GetArticleByID(ctx, id)

	// Check if the article is in the cache
	if err != article.ErrArticleCachingNotFound {
		if err != nil {
			return nil, err
		}
		return articleCache, nil
	}

	// Get the article from the database
	articleDb, err := s.articleQueryRepository.GetArticleByID(id)
	if err != nil {
		// check if the article is not found
		// gorm returns an error if the article is not found
		if err == gorm.ErrRecordNotFound {
			return nil, article.ErrArticleNotFound
		}

		return nil, err
	}

	// Create cache for the article asynchronously
	go func() {
		err := s.articleCachingRepository.CreateArticle(ctx, articleDb)
		if err != nil {
			log.Printf("failed to create cache for article: %v", err)
		}
	}()

	// Return the article
	return articleDb, nil
}

// GetListArticles retrieves a list of articles based on the given query.
//
// query: The article query parameters.
// Returns a slice of article pointers and an error if any.
func (s *ArticleService) GetListArticles(ctx context.Context, query *article.ArticleQuery) (*article.ListArticleDTO, error) {
	return s.articleQueryRepository.GetListArticles(ctx, query)
}
