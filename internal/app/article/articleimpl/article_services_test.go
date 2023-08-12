package articleimpl_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undercode99/article_service/internal/app/article"
	"github.com/undercode99/article_service/internal/app/article/articleimpl"
	"gorm.io/gorm"
)

// Mocking ArticleCommandRepository
type MockArticleCommandRepository struct {
	mock.Mock
}

func (m *MockArticleCommandRepository) CreateArticle(article *article.Article) error {
	return m.Called(article).Error(0)
}

func (m *MockArticleCommandRepository) CreateIndexArticle(ctx context.Context, item *article.Article) error {
	return m.Called(ctx, item).Error(0)
}

// Mocking ArticleQueryRepository
type MockArticleQueryRepository struct {
	mock.Mock
}

func (m *MockArticleQueryRepository) GetArticleByID(id int) (*article.Article, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.Article), args.Error(1)
}

func (m *MockArticleQueryRepository) GetListArticles(ctx context.Context, query *article.ArticleQuery) (*article.ListArticleDTO, error) {
	args := m.Called(query)
	return args.Get(0).(*article.ListArticleDTO), args.Error(1)
}

// Mocking ArticleCachingRepository
type MockArticleCachingRepository struct {
	mock.Mock
}

func (m *MockArticleCachingRepository) CreateArticle(ctx context.Context, article *article.Article) error {
	return m.Called(ctx, article).Error(0)
}

func (m *MockArticleCachingRepository) GetArticleByID(ctx context.Context, id int) (*article.Article, error) {
	args := m.Called(ctx, id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.Article), args.Error(1)
}

func TestNewArticleService(t *testing.T) {
	mockArticleCommandRepository := &MockArticleCommandRepository{}
	mockArticleQueryRepository := &MockArticleQueryRepository{}
	mockArticleCachingRepository := &MockArticleCachingRepository{}

	articleService := articleimpl.NewArticleService(mockArticleCommandRepository, mockArticleQueryRepository, mockArticleCachingRepository)

	// Testing that the returned ArticleService is not nil
	if articleService == nil {
		t.Errorf("NewArticleService() returned nil, expected non-nil ArticleService")
	}

}

// TestCreateArticle tests the CreateArticle function.
//
// This function creates a test article with a valid article create command.
// It mocks the article command repository, article query repository, and article caching repository.
// The article is created using the article service's CreateArticle method.
// The function checks that no error is returned and that the created article is not nil.
func TestCreateArticle(t *testing.T) {

	ctx := context.Background()
	// Create a valid article create command
	cmd := &article.ArticleCreateCommand{
		Title:  "Test Article",
		Body:   "This is a test article.",
		Author: "John Doe",
	}
	// Create mock repositories
	mockArticleCommandRepository := &MockArticleCommandRepository{}
	mockArticleQueryRepository := &MockArticleQueryRepository{}
	mockArticleCachingRepository := &MockArticleCachingRepository{}

	// Create an instance of the article service
	articleService := articleimpl.NewArticleService(
		mockArticleCommandRepository,
		mockArticleQueryRepository,
		mockArticleCachingRepository,
	)

	// Set up expectations for the mock repositories
	mockArticleCommandRepository.On("CreateArticle", mock.Anything).Return(nil)
	mockArticleCommandRepository.On("CreateIndexArticle", ctx, mock.Anything).Return(nil)

	// Create the article using the article service's CreateArticle method
	createdArticle, err := articleService.CreateArticle(ctx, cmd)

	// Check that no error is returned and that the created article is not nil
	assert.Nil(t, err)
	assert.NotNil(t, createdArticle)
}

// TestArticleService_GetArticleByID tests the GetArticleByID method of the ArticleService struct.
//
// 1. Test the article is found in cache.
// 2. Test the article is not found in cache.
// 3. Test the article is found in the database.
func TestArticleService_GetArticleByID(t *testing.T) {
	ctx := context.TODO()

	mockArticleCommandRepo := &MockArticleCommandRepository{}
	mockArticleQueryRepo := &MockArticleQueryRepository{}
	mockArticleCachingRepo := &MockArticleCachingRepository{}

	articleService := articleimpl.NewArticleService(mockArticleCommandRepo, mockArticleQueryRepo, mockArticleCachingRepo)

	t.Run("Article found in cache", func(t *testing.T) {
		// Mock the GetArticleByID method of the articleCachingRepository to return a non-nil article
		mockArticleCachingRepo.On("GetArticleByID", ctx, 1).Return(&article.Article{ID: 1, Title: "Test Article"}, nil)

		// Call the GetArticleByID method
		article, err := articleService.GetArticleByID(ctx, 1)

		// Assert that the article is not nil and the error is nil
		assert.NotNil(t, article)
		assert.Nil(t, err)

		// Assert that the article is the one returned by the cache
		assert.Equal(t, 1, article.ID)
		assert.Equal(t, "Test Article", article.Title)

		// Assert that the GetArticleByID method of the articleCachingRepository was called with the correct parameters
		mockArticleCachingRepo.AssertCalled(t, "GetArticleByID", ctx, 1)
	})

	t.Run("Article not found in cache", func(t *testing.T) {
		// Mock the GetArticleByID method of the articleCachingRepository to return nil
		mockArticleCachingRepo.On("GetArticleByID", ctx, 2).Return(nil, article.ErrArticleCachingNotFound)

		// Mock the GetArticleByID method of the articleQueryRepository to return an article
		mockArticleQueryRepo.On("GetArticleByID", 2).Return(&article.Article{ID: 2, Title: "Test Article 2"}, nil)

		// Mock the CreateArticle method of the articleCachingRepository
		mockArticleCachingRepo.On("CreateArticle", ctx, &article.Article{ID: 2, Title: "Test Article 2"}).Return(nil)

		// Call the GetArticleByID method
		itemsArticle, err := articleService.GetArticleByID(ctx, 2)

		// Assert that the article is not nil and the error is nil
		assert.NotNil(t, itemsArticle)
		assert.Nil(t, err)

		// Assert that the article is the one returned by the database
		assert.Equal(t, 2, itemsArticle.ID)
		assert.Equal(t, "Test Article 2", itemsArticle.Title)

		// Assert that the GetArticleByID method of the articleCachingRepository was called with the correct parameters
		mockArticleCachingRepo.AssertCalled(t, "GetArticleByID", ctx, 2)

		// Assert that the GetArticleByID method of the articleQueryRepository was called with the correct parameters
		mockArticleQueryRepo.AssertCalled(t, "GetArticleByID", 2)

		// Assert that the CreateArticle method of the articleCachingRepository was called with the correct parameters
		// mockArticleCachingRepo.AssertCalled(t, "CreateArticle", ctx, &article.Article{ID: 2, Title: "Test Article 2"})
	})

	t.Run("Article not found in database", func(t *testing.T) {
		// Mock the GetArticleByID method of the articleCachingRepository to return nil
		mockArticleCachingRepo.On("GetArticleByID", ctx, 3).Return(nil, article.ErrArticleCachingNotFound)

		// Mock the GetArticleByID method of the articleQueryRepository to return an error
		mockArticleQueryRepo.On("GetArticleByID", 3).Return(nil, gorm.ErrRecordNotFound)

		// Call the GetArticleByID method
		itemsArticle, err := articleService.GetArticleByID(ctx, 3)

		// Assert that the article is nil and the error is ErrArticleNotFound
		assert.Nil(t, itemsArticle)
		assert.Equal(t, article.ErrArticleCachingNotFound, err)

		// Assert that the GetArticleByID method of the articleCachingRepository was called with the correct parameters
		mockArticleCachingRepo.AssertCalled(t, "GetArticleByID", ctx, 3)

		// Assert that the GetArticleByID method of the articleQueryRepository was called with the correct parameters
		mockArticleQueryRepo.AssertCalled(t, "GetArticleByID", 3)
	})
}

// TestGetListArticles is a test function for the GetListArticles method of the ArticleService.
//
// It tests various scenarios of querying the article repository and asserts the returned results and errors.
// The test cases cover scenarios such as an empty query, a query with a search term, and a query with a category.
// For each test case, the GetListArticles method is called with the corresponding query, and the result and error are checked against the expected values.
// This function is used to ensure the correctness of the GetListArticles method of the ArticleService.
// It is part of the testing suite for the ArticleService.
func TestGetListArticles(t *testing.T) {
	ctx := context.Background()

	// Create mock repositories
	mockArticleCommandRepo := &MockArticleCommandRepository{}
	mockArticleQueryRepo := &MockArticleQueryRepository{}
	mockArticleCachingRepo := &MockArticleCachingRepository{}

	// Define test cases
	tests := []struct {
		name   string
		query  *article.ArticleQuery
		result *article.ListArticleDTO
		err    error
	}{
		{
			name:  "Empty query",
			query: &article.ArticleQuery{},
			result: &article.ListArticleDTO{
				Articles: []article.Article{
					{ID: 1, Title: "Test Golang", Author: "cena"},
					{ID: 2, Title: "Test Golang 2", Author: "cena"},
				},
			},
			err: nil,
		},
		{
			name:  "Query with search term",
			query: &article.ArticleQuery{Search: "golang"},
			result: &article.ListArticleDTO{
				Articles: []article.Article{
					{ID: 1, Title: "Test Golang 1", Author: "cena"},
				},
			},
			err: nil,
		},
		{
			name:  "Query with category",
			query: &article.ArticleQuery{Author: "jhon"},
			result: &article.ListArticleDTO{
				Articles: []article.Article{
					{ID: 1, Title: "Data Science", Author: "jhon"},
				},
			},
			err: nil,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new instance of the ArticleService
			articleService := articleimpl.NewArticleService(mockArticleCommandRepo, mockArticleQueryRepo, mockArticleCachingRepo)

			// Mock the GetListArticles method of the articleQueryRepository to return a non-nil article
			mockArticleQueryRepo.On("GetListArticles", tt.query).Return(tt.result, tt.err)

			// Call the GetListArticles method and retrieve the result and error
			result, err := articleService.GetListArticles(ctx, tt.query)

			assert.EqualValues(t, result.Articles, tt.result.Articles)
			assert.Equal(t, tt.err, err)
		})
	}
}
