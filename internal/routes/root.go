package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"ozon-fintech/internal/service"
)

type Router struct {
	srv     *http.Server
	logger  *slog.Logger
	storage *service.Service
}

func New(srvAddr, envLevel string, logger *slog.Logger, service *service.Service) *Router {
	switch envLevel {
	case "local":
		gin.SetMode(gin.DebugMode)
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}

	rtr := gin.Default()
	registerShortener(rtr, service)

	router := &Router{
		srv:     &http.Server{Addr: ":" + srvAddr, Handler: rtr},
		logger:  logger,
		storage: service,
	}

	return router
}

func (r *Router) Start() error {
	err := r.srv.ListenAndServe()
	return err
}

func (r *Router) Shutdown(ctx context.Context) error {
	err := r.srv.Shutdown(ctx)
	return err
}
