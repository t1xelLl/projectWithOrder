package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/projectWithOrder/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	order := router.Group("/order")
	{
		order.GET("/:order_uid", h.GetOrderByUID)
		order.POST("/", h.CreateOrder)
	}
	return router
}
