package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/undercode99/article_service/internal/app/article"
)

func (h *ApiHandler) CreateArticle(c *gin.Context) {
	var createCmd article.ArticleCreateCommand

	if err := c.BindJSON(&createCmd); err != nil {
		h.withResponseError(c, err)
		return
	}

	if err := createCmd.Validate(); err != nil {
		h.withResponseErrorStatus(c, err, http.StatusBadRequest)
		return
	}

	createdArticle, err := h.articleService.CreateArticle(c, &createCmd)
	if err != nil {
		h.withResponseError(c, err)
		return
	}

	h.withResponse(c, createdArticle, http.StatusCreated)
}

func (h *ApiHandler) GetArticleByID(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleItem, err := h.articleService.GetArticleByID(c, idInt)
	if err != nil {
		status := http.StatusInternalServerError
		if err == article.ErrArticleNotFound {
			status = http.StatusNotFound
		}
		h.withResponseErrorStatus(c, err, status)
		return
	}

	h.withResponse(c, articleItem)
}

func (h *ApiHandler) GetListArticles(c *gin.Context) {
	var qry article.ArticleQuery
	if err := c.BindQuery(&qry); err != nil {
		h.withResponseError(c, err)
		return
	}

	articles, err := h.articleService.GetListArticles(c, &qry)
	if err != nil {
		h.withResponseError(c, err)
		return
	}

	h.withResponse(c, articles)
}
