package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"github.com/t1xelLl/projectWithOrder/internal/repository"
	"github.com/t1xelLl/projectWithOrder/internal/service/cache"
	"time"
)

type OrderService struct {
	orderRepo repository.Order
	cache     cache.Cache
}

func NewOrderService(orderRepo repository.Order, cache cache.Cache) *OrderService {
	return &OrderService{orderRepo: orderRepo, cache: cache}
}

func (s *OrderService) GetOrderByUID(ctx context.Context, uid string) (*entities.Order, error) {
	cachedOrder, err := s.cache.GetData(ctx, uid)
	if err != nil && err != redis.Nil {
		logrus.Warnf("Failed to get order from cache: %v", err)
	}

	if cachedOrder != nil {
		var order entities.Order
		if err := json.Unmarshal(cachedOrder, &order); err != nil {
			return nil, fmt.Errorf("unmarshal cached order: %w", err)
		}
		return &order, err
	}
	order, err := s.orderRepo.GetOrderByUID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("get from repository %w", err)
	}

	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	go func(order *entities.Order, uid string) {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := json.Marshal(order)
		if err != nil {
			logrus.Warnf("failed to marshal order for caching %q: %v", uid, err)
			return
		}

		if err := s.cache.SetData(cacheCtx, order.OrderUID, data); err != nil {
			logrus.Warnf("failed to cache order %q: %v", order.OrderUID, err)
		} else {
			logrus.Debugf("successfully cached order %q", order.OrderUID)
		}
	}(order, uid)

	return order, nil

}

func (s *OrderService) CreateOrder(ctx context.Context, order *entities.Order) error {
	exists, err := s.orderRepo.OrderExist(ctx, order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if exists {
		return fmt.Errorf("order with UID %s already exists", order.OrderUID)
	}

	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return fmt.Errorf("create order in repository: %w", err)
	}

	go func(order *entities.Order) {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := json.Marshal(order)
		if err != nil {
			logrus.Warnf("failed to marshal order for caching %q: %v", order.OrderUID, err)
			return
		}

		if err := s.cache.SetData(cacheCtx, order.OrderUID, data); err != nil {
			logrus.Warnf("failed to cache order %q: %v", order.OrderUID, err)
		} else {
			logrus.Debugf("successfully cached order %q asynchronously", order.OrderUID)
		}
	}(order)

	return nil
}

func (s *OrderService) PreloadCache(ctx context.Context) error {
	orders, err := s.orderRepo.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to preload cache: %w", err)
	}

	for _, order := range orders {
		data, err := json.Marshal(order)
		if err != nil {
			logrus.Warnf("failed to marshal order %q: %v", order.OrderUID, err)
			continue
		}

		if err := s.cache.SetData(ctx, order.OrderUID, data); err != nil {
			logrus.Warnf("failed to cache order %q: %v", order.OrderUID, err)
		}
	}

	logrus.Infof("Preload %d orders into cache", len(orders))
	return nil
}
