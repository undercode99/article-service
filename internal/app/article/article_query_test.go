package article_test

import (
	"testing"

	"github.com/undercode99/article_service/internal/app/article"

	// assert package
	"github.com/stretchr/testify/assert"
)

// TestGetLimit is a unit test for the GetLimit function.
//
// It tests two cases:
//  1. When the limit is not set, the default value should be returned.
//  2. When the limit is set, the set value should be returned.
func TestGetLimit(t *testing.T) {
	// Test case 1: Limit is not set, default value should be returned
	a := &article.ArticleQuery{}
	expected := 10
	assert.Equal(t, expected, a.GetLimit())

	// Test case 2: Limit is set, should return the set value
	a = &article.ArticleQuery{
		Limit: 20,
	}
	expected = 20
	assert.Equal(t, expected, a.GetLimit())
}

// TestGetPage is a unit test for the GetPage function.
//
// Test cases:
//
// - Case 1: Page value is not set, should return default value of 1.
// - Case 2: Page value is set to a number value, should return the same value.
//
// This test function is meant to be used with the Go testing framework.
func TestGetPage(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		query    *article.ArticleQuery
		expected int
	}{
		{
			name:     "Page value is not set, should return default value of 1",
			query:    &article.ArticleQuery{},
			expected: 1,
		},
		{
			name:     "Page value is set to a number value, should return the same value",
			query:    &article.ArticleQuery{Page: 5},
			expected: 5,
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the GetPage function on the query
			result := tc.query.GetPage()
			// Check if the result matches the expected value
			assert.Equal(t, tc.expected, result)
		})
	}
}
