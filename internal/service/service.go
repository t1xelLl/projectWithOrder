package service

import (
	"context"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"github.com/t1xelLl/projectWithOrder/internal/repository"
	"github.com/t1xelLl/projectWithOrder/internal/service/cache"
)

type Order interface {
	GetOrderByUID(ctx context.Context, uid string) (*entities.Order, error)
	CreateOrder(ctx context.Context, order *entities.Order) error
	PreloadCache(ctx context.Context) error
}

type Service struct {
	Order
}

func NewService(repo *repository.Repository, cache *cache.Cache) *Service {
	return &Service{
		Order: NewOrderService(repo, *cache),
	}
}
