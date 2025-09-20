package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

// TODO: CHANGE
func (r *OrderPostgres) GetOrderByUID(uid string) (entities.Order, error) {
	var order entities.Order

	query := `
		SELECT 
			order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service, 
			shardkey, sm_id, date_created, oof_shard
		FROM "order" 
		WHERE order_uid = $1
	`

	err := r.db.Get(&order, query, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.Order{}, fmt.Errorf("order not found")
		}
		return entities.Order{}, fmt.Errorf("failed to get order : %w", err)
	}

	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery 
		WHERE order_uid = $1
	`
	err = r.db.Get(&order.Delivery, deliveryQuery, uid)
	if err != nil {
		return entities.Order{}, fmt.Errorf("failed to get delivery info: %w", err)
	}

	paymentQuery := `
		SELECT transaction, request_id, currency, provider, amount, 
		       payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment 
		WHERE order_uid = $1
	`
	err = r.db.Get(&order.Payment, paymentQuery, uid)
	if err != nil {
		return entities.Order{}, fmt.Errorf("failed to get payment info: %w", err)
	}

	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, sale, 
		       size, total_price, nm_id, brand, status
		FROM item 
		WHERE order_uid = $1
	`
	err = r.db.Select(&order.Items, itemsQuery, uid)
	if err != nil {
		return entities.Order{}, fmt.Errorf("failed to get items: %w", err)
	}

	return order, nil

}

func (r *OrderPostgres) CreateOrder(order entities.Order) error {
	// Начинаем транзакцию
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Откатываем в случае ошибки

	// Вставляем основную информацию о заказе
	orderQuery := `
		INSERT INTO "order" (
			order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service, 
			shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(orderQuery,
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

	// Вставляем информацию о доставке
	deliveryQuery := `
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.Exec(deliveryQuery,
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

	// Вставляем информацию о платеже
	paymentQuery := `
		INSERT INTO payment (
			order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.Exec(paymentQuery,
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

	// Вставляем товары
	itemQuery := `
		INSERT INTO item (
			order_uid, chrt_id, track_number, price, rid, name, 
			sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	for _, item := range order.Items {
		_, err = tx.Exec(itemQuery,
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

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
