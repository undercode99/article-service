package article

import (
	"context"
	"errors"
	"time"
)

var (
	ErrArticleNotFound        = errors.New("article not found")
	ErrArticleValidation      = errors.New("article validation error")
	ErrArticleCachingNotFound = errors.New("article not found")

	// IndexName is the name of the index in Elasticsearch
	IndexName = "articles"
)

type Article struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
}

// NewArticle creates a new article based on the provided ArticleCreateCommand.
//
// The function takes a pointer to an ArticleCreateCommand as its parameter
// and returns a pointer to an Article. The Article struct is populated with
// the values from the command parameter, including the author, title, body,
// and creation timestamp.
func NewArticle(cmd *ArticleCreateCommand) *Article {
	return &Article{
		Author:  cmd.Author,
		Title:   cmd.Title,
		Body:    cmd.Body,
		Created: time.Now(),
	}
}

type ArticleCommandRepository interface {
	CreateArticle(article *Article) error
	CreateIndexArticle(ctx context.Context, article *Article) error
}

type ArticleCachingRepository interface {
	CreateArticle(ctx context.Context, article *Article) error
	GetArticleByID(ctx context.Context, id int) (*Article, error)
}

type ArticleQueryRepository interface {
	GetArticleByID(id int) (*Article, error)
	GetListArticles(ctx context.Context, query *ArticleQuery) (*ListArticleDTO, error)
}

type ArticleService interface {
	CreateArticle(ctx context.Context, cmd *ArticleCreateCommand) (*Article, error)
	GetArticleByID(ctx context.Context, id int) (*Article, error)
	GetListArticles(ctx context.Context, query *ArticleQuery) (*ListArticleDTO, error)
}
