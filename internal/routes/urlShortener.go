package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"ozon-fintech/internal/models"
	"ozon-fintech/internal/service"

	"github.com/gin-gonic/gin"
)

type implShortener struct {
	r       *gin.Engine
	service *service.Service
}

func registerShortener(r *gin.Engine, service *service.Service) {
	shortenerImpl := implShortener{r: r, service: service}

	r.GET("/getFullURL", shortenerImpl.getFullURL)
	r.POST("/loadFullURL", shortenerImpl.loadShortURL)
}

func (impl *implShortener) getFullURL(ctx *gin.Context) {
	shortURL := ctx.Query("short_url")
	err := service.ValidShortURL(shortURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	fullURL, err := impl.service.GetFullURL(shortURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, err)
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}
	ctx.JSON(http.StatusOK, fullURL)
}

func (impl *implShortener) loadShortURL(ctx *gin.Context) {
	fullURL := ctx.Query("full_url")
	err := service.ValidFullURL(fullURL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	link := models.Link{FullUrl: fullURL}
	shortURL, err := impl.service.LoadShortURL(link)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, shortURL)
}
