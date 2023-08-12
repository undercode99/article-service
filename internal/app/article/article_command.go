package article

import "errors"

var (
	ErrAuthorIsRequired = errors.New("author is required")
	ErrTitleIsRequired  = errors.New("title is required")
	ErrBodyIsRequired   = errors.New("body is required")
)

type ArticleCreateCommand struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Validate checks if the ArticleCreateCommand is valid.
//
// It returns an error if any of the required fields are missing.
func (a *ArticleCreateCommand) Validate() error {
	if a.Author == "" {
		return ErrAuthorIsRequired
	}
	if a.Title == "" {
		return ErrTitleIsRequired
	}
	if a.Body == "" {
		return ErrBodyIsRequired
	}

	return nil
}
