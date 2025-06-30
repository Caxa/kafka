// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"
	"time"

	"l0/cmd/internal/models"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func NewPostgres(dsn string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	db, err = sql.Open("postgres", dsn)

	for i := 0; i < 3 && err != nil; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			time.Sleep(5 * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}

	Conn = db
	return db, nil
}

func InsertOrder(db *sql.DB, order models.Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1,$2,$3,$4,'',$5,$6,$7,$8,$9,$10)
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.CustomerID, order.DeliveryService, order.ShardKey,
		order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO deliveries (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("insert delivery: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO payments (
			order_uid, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	for _, item := range order.Items {
		_, err := tx.Exec(`
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, sale,
				size, total_price, nm_id, brand, status
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
	}

	return tx.Commit()
}

func GetOrderByID(db *sql.DB, id string) (models.Order, error) {
	var order models.Order

	err := db.QueryRow(`
		SELECT order_uid, track_number, entry, locale, customer_id,
		       delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
	`, id).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry,
		&order.Locale, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return order, err
	}

	err = db.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries WHERE order_uid = $1
	`, id).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return order, err
	}

	err = db.QueryRow(`
		SELECT transaction, request_id, currency, provider,
		       amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payments WHERE order_uid = $1
	`, id).Scan(&order.Payment.Transaction, &order.Payment.RequestID,
		&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
		&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return order, err
	}

	rows, err := db.Query(`
		SELECT chrt_id, track_number, price, rid, name, sale,
		       size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
	`, id)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func GetAllOrders(db *sql.DB) ([]models.Order, error) {
	rows, err := db.Query(`SELECT order_uid FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		order, err := GetOrderByID(db, id)
		if err == nil {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
