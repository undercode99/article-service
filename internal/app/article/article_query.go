package article

type ArticleQuery struct {
	Search     string `form:"search"`
	Author     string `form:"author"`
	SortNewest bool   `form:"sort_newest"`
	Limit      int    `form:"limit"`
	Page       int    `form:"page"`
}

// GetLimit returns the limit value of the ArticleQuery for pagination purposes.
//
// It does not accept any parameters.
// It returns an int value.
func (a *ArticleQuery) GetLimit() int {
	// Set the limit value to the default value of 10 if it is not set.
	if a.Limit == 0 {
		a.Limit = 10
	}

	return a.Limit
}

// GetPage returns the page number for the ArticleQuery for pagination purposes.
//
// It does not take any parameters.
// It returns an integer.
func (a *ArticleQuery) GetPage() int {
	// If the page value is not set, set it to the default value of 1.
	if a.Page == 0 {
		return 1
	}
	return a.Page
}
