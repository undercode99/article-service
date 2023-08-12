package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/undercode99/article_service/internal/api"
	"github.com/undercode99/article_service/internal/app/article"

	"testing"
)

type mockArticleService struct {
}

func (m *mockArticleService) CreateArticle(ctx context.Context, cmd *article.ArticleCreateCommand) (*article.Article, error) {
	return &article.Article{
		ID:     1,
		Title:  cmd.Title,
		Body:   cmd.Body,
		Author: cmd.Author,
	}, nil
}

func (m *mockArticleService) GetArticleByID(ctx context.Context, id int) (*article.Article, error) {
	if id != 1 {
		return nil, article.ErrArticleNotFound
	}

	return &article.Article{
		ID:     1,
		Title:  "Test Article",
		Body:   "This is a test article",
		Author: "Lorem ipsum dolor sit amet",
	}, nil
}

func (m *mockArticleService) GetListArticles(ctx context.Context, query *article.ArticleQuery) (*article.ListArticleDTO, error) {
	return &article.ListArticleDTO{
		Articles: []article.Article{
			{
				ID:     1,
				Title:  "Test Article",
				Body:   "This is a test article",
				Author: "Lorem ipsum dolor sit amet",
			},
		},
		Page:  1,
		Limit: 10,
	}, nil
}

func TestApiHandler_CreateArticle(t *testing.T) {
	// Test case: successful creation of an article
	t.Run("Successful creation", func(t *testing.T) {
		mockArticleService := &mockArticleService{}

		// Create the API handler with the mock article service
		apiHandler := api.NewApiHandler(mockArticleService)

		r := gin.Default()
		r.POST("/v1/articles", apiHandler.CreateArticle)

		payload := `{"title": "Test Article", "body": "This is a test article", "author": "jhon"}`
		req, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		bodyResultMap := make(map[string]interface{})
		json.Unmarshal(w.Body.Bytes(), &bodyResultMap)

		assert.Equal(t, "Test Article", bodyResultMap["title"])
		assert.Equal(t, "This is a test article", bodyResultMap["body"])
		assert.Equal(t, "jhon", bodyResultMap["author"])
	})

	t.Run("Unsuccessful creation", func(t *testing.T) {
		mockArticleService := &mockArticleService{}

		// Create the API handler with the mock article service
		apiHandler := api.NewApiHandler(mockArticleService)

		r := gin.Default()
		r.POST("/v1/articles", apiHandler.CreateArticle)

		payload := `{"title": "Test Article", "body": "This is a test article"}`
		req, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestApiHandler_GetArticleByID(t *testing.T) {
	// Test case: successful retrieval of an article
	t.Run("Successful retrieval", func(t *testing.T) {
		mockArticleService := &mockArticleService{}

		// Create the API handler with the mock article service
		apiHandler := api.NewApiHandler(mockArticleService)

		r := gin.Default()
		r.GET("/v1/articles/:id", apiHandler.GetArticleByID)

		req, _ := http.NewRequest("GET", "/v1/articles/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Unsuccessful retrieval", func(t *testing.T) {
		mockArticleService := &mockArticleService{}

		// Create the API handler with the mock article service
		apiHandler := api.NewApiHandler(mockArticleService)

		r := gin.Default()
		r.GET("/v1/articles/:id", apiHandler.GetArticleByID)

		req, _ := http.NewRequest("GET", "/v1/articles/2", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

	})
}
