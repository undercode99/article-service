package articleimpl

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/undercode99/article_service/internal/app/article"
)

type ArticleCachingRepository struct {
	redisClient *redis.Client
}

// NewArticleCachingRepository returns a new instance of article.ArticleCachingRepository.
//
// It takes a redisClient as a parameter and initializes the ArticleCachingRepository
// struct with the provided redisClient.
// It returns a pointer to the initialized ArticleCachingRepository.
func NewArticleCachingRepository(redisClient *redis.Client) article.ArticleCachingRepository {
	return &ArticleCachingRepository{
		redisClient: redisClient,
	}
}

// CreateArticle creates a new article in the caching repository.
//
// ctx - The context in which the function is being called.
// article - The article object to be created.
// Returns an error if there was an issue creating the article.
func (r *ArticleCachingRepository) CreateArticle(ctx context.Context, article *article.Article) error {

	// create cache key
	key := "article:" + strconv.Itoa(article.ID)

	articleJSON, err := json.Marshal(article)
	if err != nil {
		return err
	}

	// set cache
	return r.redisClient.Set(ctx, key, articleJSON, time.Hour*20).Err()
}

// GetFromCache retrieves an article from the cache based on its ID.
//
// ctx - the context.Context object used for cancellation and timeouts.
// id - the ID of the article to retrieve.
// Returns a pointer to the retrieved article and an error, if any.
func (r *ArticleCachingRepository) GetArticleByID(ctx context.Context, id int) (*article.Article, error) {
	key := "article:" + strconv.Itoa(id)
	// check if the article is in the cache
	exists, err := r.redisClient.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("failed to check if article is in cache: %v", err)
		return nil, err
	}
	if exists == 0 {
		log.Printf("article not found in cache")
		return nil, article.ErrArticleCachingNotFound
	}

	// get the article from the cache
	var article article.Article
	articleJSON, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Printf("failed to get article from cache: %v", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(articleJSON), &article)
	if err != nil {
		log.Printf("failed to unmarshal article from cache: %v", err)
		return nil, err
	}

	return &article, nil
}
