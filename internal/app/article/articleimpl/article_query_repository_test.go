package articleimpl_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
	"github.com/undercode99/article_service/internal/app/article/articleimpl"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockTransport struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func dbMockConnection() (*gorm.DB, sqlmock.Sqlmock) {
	mockDb, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{})
	return db, mock

}

func elasticMockConnection() (*elasticsearch.TypedClient, MockTransport) {
	mocktrans := MockTransport{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mocktrans.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mocktrans.Response, nil }

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Transport: &mocktrans,
	})
	if err != nil {
		panic(err)
	}
	return client, mocktrans
}

// TestGetArticleByID tests the GetArticleByID function.
//
// It checks if the article with the given ID exists in the database
// and verifies that the returned article has the correct ID and title.
// It uses a mock database connection and a mock Elasticsearch client
// for testing purposes.
//
// Parameters:
// - t: The testing.T object for running tests and reporting results.
func TestGetArticleByID(t *testing.T) {

	db, mock := dbMockConnection()
	client, _ := elasticMockConnection()
	repo := articleimpl.NewArticleQueryRepository(db, client)

	t.Run("Test exists article", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM \"articles\" WHERE \"articles\".\"id\" = ?").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Test Article"))

		res, err := repo.GetArticleByID(1)
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}

		assert.Equal(t, 1, res.ID)
		assert.Equal(t, "Test Article", res.Title)
	})

}

func TestGetListArticlesElastic(t *testing.T) {
	// TODO: Implement test cases for GetListArticlesElastic function
}
