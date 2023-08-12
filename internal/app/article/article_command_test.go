package article_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/undercode99/article_service/internal/app/article"
)

// TestArticleCreateCommand_Validate tests the validation of the ArticleCreateCommand function.
//
// It tests different scenarios by providing various command inputs and expected error outputs.
// The function iterates over a list of test cases and performs the validation for each one.
// It uses the assert.Equal function to compare the expected error with the actual error returned.
func TestArticleCreateCommand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		command *article.ArticleCreateCommand
		wantErr error
	}{
		{
			name: "valid command",
			command: &article.ArticleCreateCommand{
				Author: "John Doe",
				Title:  "Test Article",
				Body:   "This is a test article",
			},
			wantErr: nil,
		},
		{
			name: "missing author",
			command: &article.ArticleCreateCommand{
				Author: "",
				Title:  "Test Article",
				Body:   "This is a test article",
			},
			wantErr: article.ErrAuthorIsRequired,
		},
		{
			name: "missing title",
			command: &article.ArticleCreateCommand{
				Author: "John Doe",
				Title:  "",
				Body:   "This is a test article",
			},
			wantErr: article.ErrTitleIsRequired,
		},
		{
			name: "missing body",
			command: &article.ArticleCreateCommand{
				Author: "John Doe",
				Title:  "Test Article",
				Body:   "",
			},
			wantErr: article.ErrBodyIsRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.command.Validate()
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
