package api

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/undercode99/article_service/config"
)

type ApiService struct {
	apiHandler *ApiHandler
	cfg        *config.Config
}

func NewApiService(apiHandler *ApiHandler, cfg *config.Config) *ApiService {
	return &ApiService{
		apiHandler: apiHandler,
		cfg:        cfg,
	}
}

// Run runs the API service.
//
// It sets up the routes for the API endpoints and starts the server.
// The function takes a context.Context as a parameter and returns nothing.
func (a *ApiService) Run(ctx context.Context) {
	r := gin.Default()

	// set gin mode
	if a.cfg.AppModeIsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// api routes
	v1 := r.Group("/v1")
	{
		v1.POST("/articles", a.apiHandler.CreateArticle)
		v1.GET("/articles/:id", a.apiHandler.GetArticleByID)
		v1.GET("/articles", a.apiHandler.GetListArticles)
	}

	log.Printf("Starting server on port %s", a.cfg.AppPort)
	err := r.Run(":" + a.cfg.AppPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
