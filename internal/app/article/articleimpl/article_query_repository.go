package articleimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/undercode99/article_service/internal/app/article"
	"gorm.io/gorm"
)

type ArticleQueryRepository struct {
	db            *gorm.DB
	elasticClient *elasticsearch.TypedClient
}

func NewArticleQueryRepository(db *gorm.DB, elasticClient *elasticsearch.TypedClient) article.ArticleQueryRepository {
	return &ArticleQueryRepository{
		db:            db,
		elasticClient: elasticClient,
	}
}

// GetByID returns an article by its ID.
//
// It takes an integer parameter representing the ID of the article to retrieve.
// The function returns a pointer to the retrieved article and an error, if any.
func (a ArticleQueryRepository) GetArticleByID(id int) (*article.Article, error) {
	var article article.Article
	if err := a.db.First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetListArticles retrieves a list of articles based on the provided query.
//
// ctx: The context in which the function is being executed.
// qry: The article query object containing the search parameters.
// Returns a ListArticleDTO and an error.
func (a ArticleQueryRepository) GetListArticles(ctx context.Context, qry *article.ArticleQuery) (*article.ListArticleDTO, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []map[string]interface{}{},
				"should": []map[string]interface{}{},
			},
		},
		"sort": map[string]interface{}{
			"created": map[string]interface{}{
				"order": "asc",
			},
		},
		"from": (qry.GetPage() - 1) * qry.GetLimit(),
		"size": qry.GetLimit(),
	}

	if qry.Author != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = []map[string]interface{}{
			{
				"match": map[string]interface{}{
					"author": qry.Author,
				},
			},
		}
	}

	if qry.Search != "" {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = []map[string]interface{}{
			{
				"match": map[string]interface{}{
					"title": qry.Search,
				},
			},
			{
				"match": map[string]interface{}{
					"body": qry.Search,
				},
			},
		}
	}
	if qry.SortNewest {
		query["sort"] = map[string]interface{}{
			"created": map[string]interface{}{
				"order": "desc",
			},
		}
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{article.IndexName},
		Body:  &buf,
	}

	res, err := req.Do(ctx, a.elasticClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var docs struct {
		Hits struct {
			Hits []struct {
				Article article.Article `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&docs); err != nil {
		return nil, err
	}

	articles := make([]article.Article, len(docs.Hits.Hits))
	for i, hit := range docs.Hits.Hits {
		articles[i] = hit.Article
	}

	return &article.ListArticleDTO{
		Articles: articles,
		Page:     qry.GetPage(),
		Limit:    qry.GetLimit(),
	}, nil
}
