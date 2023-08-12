//go:build wireinject
// +build wireinject

package runner

import (
	"context"
	"github.com/google/wire"
	"github.com/undercode99/article_service/config"
	"github.com/undercode99/article_service/internal/api"
	"github.com/undercode99/article_service/internal/app/article/articleimpl"
	"github.com/undercode99/article_service/internal/caching"
	"github.com/undercode99/article_service/internal/database"
	"github.com/undercode99/article_service/internal/searching"
)

var appSet = wire.NewSet(
	config.NewConfig,
	caching.NewRedisCaching,
	database.NewDatabase,
	searching.NewElasticClient,
	api.NewApiHandler,
	api.NewApiService,
	NewAppRunner,
)

var repositorySet = wire.NewSet(
	appSet,
	articleimpl.NewArticleQueryRepository,
	articleimpl.NewArticleCommandRepository,
	articleimpl.NewArticleCachingRepository,
)

var serviceSet = wire.NewSet(
	repositorySet,
	articleimpl.NewArticleService,
)

func InitializeApp(ctx context.Context) *AppRunner {
	wire.Build(serviceSet)

	// return valies
	return nil
}
