package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"time"
)

const operationTimeout = 10 * time.Second

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

// TODO: add context

func (r *OrderPostgres) GetOrderByUID(ctx context.Context, uid string) (*entities.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var order entities.Order

	query := `
		SELECT 
			order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service, 
			shardkey, sm_id, date_created, oof_shard
		FROM "order" 
		WHERE order_uid = $1
	`

	err := r.db.GetContext(ctx, &order, query, uid)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order : %w", err)
	}

	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery 
		WHERE order_uid = $1
	`
	err = r.db.GetContext(ctx, &order.Delivery, deliveryQuery, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery info: %w", err)
	}

	paymentQuery := `
		SELECT transaction, request_id, currency, provider, amount, 
		       payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment 
		WHERE order_uid = $1
	`
	err = r.db.GetContext(ctx, &order.Payment, paymentQuery, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment info: %w", err)
	}

	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, sale, 
		       size, total_price, nm_id, brand, status
		FROM item 
		WHERE order_uid = $1
	`
	err = r.db.SelectContext(ctx, &order.Items, itemsQuery, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return &order, nil

}

// TODO: add context

func (r *OrderPostgres) CreateOrder(ctx context.Context, order *entities.Order) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logrus.Errorf("failed to rollback transaction: %v", rbErr)
			}
		}
	}()

	orderQuery := `
		INSERT INTO "order" (
			order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service, 
			shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.ExecContext(ctx, orderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	deliveryQuery := `
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, deliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	paymentQuery := `
		INSERT INTO payment (
			order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.ExecContext(ctx, paymentQuery,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	itemQuery := `
		INSERT INTO item (
			order_uid, chrt_id, track_number, price, rid, name, 
			sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TODO: add context

func (r *OrderPostgres) GetAllOrders(ctx context.Context) ([]*entities.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var orderUIDs []string
	query := `SELECT order_uid FROM "order" ORDER BY date_created DESC LIMIT 1000`

	err := r.db.SelectContext(ctx, &orderUIDs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get order uids: %w", err)
	}

	orders := make([]*entities.Order, 0, len(orderUIDs))
	for _, uid := range orderUIDs {
		order, err := r.GetOrderByUID(ctx, uid)
		if err != nil {
			logrus.Warnf("failed to get order %s: %v", uid, err)
			continue
		}
		orders = append(orders, order)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found")
	}

	return orders, nil
}

func (r *OrderPostgres) OrderExist(ctx context.Context, uid string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	var exist bool
	query := `SELECT EXISTS (SELECT 1 FROM "order" WHERE order_uid = $1)`

	err := r.db.GetContext(ctx, &exist, query, uid)
	if err != nil {
		return false, fmt.Errorf("failed to check if order exists: %w", err)
	}
	return exist, nil
}
