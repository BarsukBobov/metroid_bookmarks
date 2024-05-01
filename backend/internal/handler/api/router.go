package api

import (
	"github.com/gin-gonic/gin"
	"metroid_bookmarks/internal/handler/api/base_api"
	"metroid_bookmarks/internal/handler/api/v1"
	"metroid_bookmarks/internal/service"
)

type router struct {
	*baseApi.Router
	service *service.Service
}

func NewRouter(
	baseAPIRouter *baseApi.Router,
	service *service.Service,
) baseApi.ApiRouter {
	return &router{
		Router:  baseAPIRouter,
		service: service,
	}
}

func (h *router) RegisterHandlers(router *gin.RouterGroup) {
	v1Group := router.Group("/v1", h.Middleware.SessionRequired)
	v1Router := v1.NewRouter(h.Router, h.service)
	v1Router.RegisterHandlers(v1Group)
}