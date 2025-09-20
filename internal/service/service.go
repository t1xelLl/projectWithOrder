package service

import (
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"github.com/t1xelLl/projectWithOrder/internal/repository"
)

type Order interface {
	GetOrderByUID(uid string) (entities.Order, error)
	CreateOrder(order entities.Order) error
}

type Service struct {
	Order
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repo.Order),
	}
}
