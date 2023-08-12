package articleimpl_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/undercode99/article_service/internal/app/article"
	"github.com/undercode99/article_service/internal/app/article/articleimpl"
)

// TestCreateArticleDatabase tests the CreateArticle function
// which creates an article in the database.
//
// This function expects a testing.T object as a parameter and does not return anything.
func TestCreateArticleDatabase(t *testing.T) {
	// Create a mock database connection
	db, mock := dbMockConnection()

	// Create a mock ElasticSearch connection
	client, _ := elasticMockConnection()

	// Create a test article
	item := article.Article{
		Title:   "Test Article",
		Body:    "This is a test article.",
		Author:  "John Doe",
		Created: time.Now(),
	}

	// Begin the transaction
	mock.ExpectBegin()

	// Expect an insert query and return the inserted ID
	mock.ExpectQuery("INSERT INTO (.+)").WithArgs(
		item.Title,
		item.Body,
		item.Author,
		item.Created,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(item.ID))

	// Commit the transaction
	mock.ExpectCommit()

	// Create a new repository instance
	repo := articleimpl.NewArticleCommandRepository(db, client)

	// Create the article using the repository
	err := repo.CreateArticle(&item)
	assert.NoError(t, err)
}

func TestCreateIndexArticle(t *testing.T) {
	// TODO: Implement test cases for CreateIndexArticle function
}
