package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/undercode99/article_service/internal/app/article"
)

type ApiHandler struct {
	articleService article.ArticleService
}

func NewApiHandler(articleService article.ArticleService) *ApiHandler {
	return &ApiHandler{
		articleService: articleService,
	}
}

func (h *ApiHandler) withResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

func (h *ApiHandler) withResponseErrorStatus(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{"error": err.Error()})
}

func (h *ApiHandler) withResponse(c *gin.Context, data interface{}, status ...int) {
	if len(status) > 0 {
		c.JSON(status[0], data)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
