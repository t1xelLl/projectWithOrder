package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"net/http"
)

// TODO: CHANGE
func (h *Handler) GetOrderByUID(c *gin.Context) {
	orderUID := c.Param("order_uid")
	if orderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uid is required"})
		return
	}

	order, err := h.services.GetOrderByUID(orderUID)
	if err != nil {
		logrus.Errorf("get order by uid %s error: %v", orderUID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order Not Found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var order entities.Order

	if err := c.BindJSON(&order); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if order.OrderUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uid is required"})
		return
	}

	if err := h.services.Order.CreateOrder(order); err != nil {
		logrus.Errorf("Failed to create order %s: %v", order.OrderUID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Order created successfully",
		"order_uid": order.OrderUID,
	})
}
