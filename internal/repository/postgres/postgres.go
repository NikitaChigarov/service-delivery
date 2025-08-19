package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"service-delivery/internal/domain"
	"time"

	_ "github.com/lib/pq"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(dsn string) (*OrderRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &OrderRepository{db: db}, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *domain.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert into orders table
	_, err = tx.ExecContext(ctx, `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}

	// Insert delivery
	_, err = tx.ExecContext(ctx, `
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email`,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("failed to insert delivery: %v", err)
	}

	// Insert payment
	_, err = tx.ExecContext(ctx, `
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction,
			request_id = EXCLUDED.request_id,
			currency = EXCLUDED.currency,
			provider = EXCLUDED.provider,
			amount = EXCLUDED.amount,
			payment_dt = EXCLUDED.payment_dt,
			bank = EXCLUDED.bank,
			delivery_cost = EXCLUDED.delivery_cost,
			goods_total = EXCLUDED.goods_total,
			custom_fee = EXCLUDED.custom_fee`,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %v", err)
	}

	// Delete existing items and insert new ones
	_, err = tx.ExecContext(ctx, "DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return fmt.Errorf("failed to delete existing items: %v", err)
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, 
				sale, size, total_price, nm_id, brand, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("failed to insert item: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *OrderRepository) GetOrderByUID(ctx context.Context, orderUID string) (*domain.Order, error) {
	var order domain.Order
	var dateCreated string

	// Get order info
	err := r.db.QueryRowContext(ctx, `
		SELECT 
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1`, orderUID).
		Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &dateCreated, &order.OofShard)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	order.DateCreated = dateCreated

	// Get delivery info
	err = r.db.QueryRowContext(ctx, `
		SELECT name, phone, zip, city, address, region, email 
		FROM deliveries 
		WHERE order_uid = $1`, orderUID).
		Scan(
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %v", err)
	}

	// Get payment info
	err = r.db.QueryRowContext(ctx, `
		SELECT 
			transaction, request_id, currency, provider, amount, 
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments 
		WHERE order_uid = $1`, orderUID).
		Scan(
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %v", err)
	}

	// Get items
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			chrt_id, track_number, price, rid, name, sale, 
			size, total_price, nm_id, brand, status
		FROM items 
		WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %v", err)
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning items: %v", err)
	}

	return &order, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT order_uid FROM orders")
	if err != nil {
		return nil, fmt.Errorf("failed to get order UIDs: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return nil, fmt.Errorf("failed to scan order UID: %v", err)
		}

		order, err := r.GetOrderByUID(ctx, orderUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get order %s: %v", orderUID, err)
		}
		if order != nil {
			orders = append(orders, *order)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning order UIDs: %v", err)
	}

	return orders, nil
}
