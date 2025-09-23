package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
)

type Order interface {
	GetOrderByUID(ctx context.Context, uid string) (*entities.Order, error)
	CreateOrder(ctx context.Context, order *entities.Order) error
	GetAllOrders(ctx context.Context) ([]*entities.Order, error)
	OrderExist(ctx context.Context, uid string) (bool, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Order: NewOrderPostgres(db)}
}
