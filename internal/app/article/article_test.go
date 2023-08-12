package article_test

import (
	"testing"

	"github.com/undercode99/article_service/internal/app/article"

	"github.com/stretchr/testify/assert"
)

// TestNewArticle tests the NewArticle function.
//
// It creates a new article with valid command values and verifies that the
// author, title, and body of the new article match the values provided in the
// command.
// It also creates a new article with empty command values and verifies that the
// author, title, and body of the new article are empty.
func TestNewArticle(t *testing.T) {
	// Test case 1: Create a new article with valid command values
	cmd := &article.ArticleCreateCommand{
		Author: "John Doe",
		Title:  "Test Article",
		Body:   "This is a test article.",
	}
	newArticle := article.NewArticle(cmd)
	assert.Equal(t, cmd.Author, newArticle.Author, "Expected author to be %s", cmd.Author)
	assert.Equal(t, cmd.Title, newArticle.Title, "Expected title to be %s", cmd.Title)
	assert.Equal(t, cmd.Body, newArticle.Body, "Expected body to be %s", cmd.Body)

	// Test case 2: Create a new article with empty command values
	cmd = &article.ArticleCreateCommand{}
	newArticle = article.NewArticle(cmd)
	assert.Empty(t, newArticle.Author, "Expected author to be empty")
	assert.Empty(t, newArticle.Title, "Expected title to be empty")
	assert.Empty(t, newArticle.Body, "Expected body to be empty")
}
