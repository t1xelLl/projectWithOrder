package service

import (
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"github.com/t1xelLl/projectWithOrder/internal/repository"
)

type OrderService struct {
	orderRepo repository.Order
}

func NewOrderService(orderRepo repository.Order) *OrderService {
	return &OrderService{orderRepo}
}

// TODO: CHANGE
func (s *OrderService) GetOrderByUID(uid string) (entities.Order, error) {
	return s.orderRepo.GetOrderByUID(uid)
}

func (s *OrderService) CreateOrder(order entities.Order) error {
	return s.orderRepo.CreateOrder(order)
}
