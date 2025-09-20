package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
)

type Order interface {
	GetOrderByUID(uid string) (entities.Order, error)
	CreateOrder(order entities.Order) error
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{Order: NewOrderPostgres(db)}
}
