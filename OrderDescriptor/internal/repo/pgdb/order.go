package pgdb

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
)

type OrderRepo struct {
	pg *postgres.Postgres
}

func NewOrderRepo(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{
		pg: pg,
	}
}

func (r *OrderRepo) CreateOrder(order model.Order) error {
	sql := `
  	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
  	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
 	`

	_, err := r.pg.Pool.Exec(context.Background(), sql,
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
		return fmt.Errorf("failed to insert new order: %v", err)
	}
	return nil
}

func (r *OrderRepo) GetOrder(uid string) (*model.Order, error) {
	order := &model.Order{}

	query := `SELECT * FROM orders WHERE order_uid = $1`
	err := r.pg.Pool.QueryRow(context.Background(), query, uid).
		Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)

	if err != nil {
		return nil, fmt.Errorf("failed to get order: %v", err)
	}
	return order, nil
}

func (r *OrderRepo) GetAllOrders() []model.Order {
	var orders []model.Order

	query := `SELECT * FROM orders`
	rows, err := r.pg.Pool.Query(context.Background(), query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			log.Errorf("failed to scan orders: %v", err)
			return nil
		}

		orders = append(orders, order)
	}

	log.Infoln("Orders wrote to memory")
	return orders
}

func (r *OrderRepo) DeleteOrder(uid string) error {
	sql := `DELETE FROM orders WHERE order_uid = $1`
	_, err := r.pg.Pool.Exec(context.Background(), sql, uid)

	if err != nil {
		log.Errorf("failed to delete order from db: %v", err)
		return err
	}
	log.Infof("Order %s deleted from db", uid)
	return nil
}

func (r *OrderRepo) DeleteAllOrders() error {
	sql := `DELETE FROM orders`
	_, err := r.pg.Pool.Exec(context.Background(), sql)

	if err != nil {
		log.Errorf("failed to delete all orders from db: %v", err)
		return err
	}
	log.Infoln("All orders deleted from db")
	return nil
}
