package routes

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Router struct {
	router *gin.Engine
	logger *slog.Logger
}

func New(rtr *gin.Engine, logger *slog.Logger) *Router {
	return &Router{
		router: rtr,
		logger: logger,
	}
}

func (r *Router) Start(addr string) error {
	err := r.router.Run(addr)
	if err != nil {
		return err
	}
	return nil
}
