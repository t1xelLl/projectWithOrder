package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/t1xelLl/projectWithOrder/internal/service"
)

type Handler struct {
	services *service.Service
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(CORSMiddleware())
	router.StaticFile("/", "./web/index.html")
	order := router.Group("/order")
	{
		order.GET("/:order_uid", h.GetOrderByUID)
		order.POST("/", h.CreateOrder)
	}
	return router
}
