package pgdb

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
)

type DeliveryRepo struct {
	pg *postgres.Postgres
}

func NewDeliveryRepo(pg *postgres.Postgres) *DeliveryRepo {
	return &DeliveryRepo{pg}
}

func (r *DeliveryRepo) CreateDelivery(delivery model.Delivery) error {
	sql := `
  	INSERT INTO deliveries (client_name, phone, zip, city, address, region, email, order_id)
  	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
 	`

	_, err := r.pg.Pool.Exec(context.Background(), sql,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
		delivery.OrderID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *DeliveryRepo) GetDelivery(id string) (*model.Delivery, error) {
	delivery := &model.Delivery{}

	query := `SELECT * FROM deliveries WHERE order_id = $1`
	err := r.pg.Pool.QueryRow(context.Background(), query, id).
		Scan(
			&delivery.ID,
			&delivery.OrderID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %v", err)
	}

	return delivery, nil
}

func (r *DeliveryRepo) GetAllDeliveries() []model.Delivery {
	deliveries := make([]model.Delivery, 0)

	query := `SELECT * FROM deliveries`
	rows, err := r.pg.Pool.Query(context.Background(), query)
	if err != nil {
		log.Errorf("Failed to get deliveries: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		delivery := model.Delivery{}
		err := rows.Scan(
			&delivery.ID,
			&delivery.OrderID,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
		)
		if err != nil {
			log.Errorf("Failed to write delivery to memory: %v", err)
			return nil
		}
		deliveries = append(deliveries, delivery)
	}

	log.Infoln("Deliveries wrote to memory")
	return deliveries
}

func (r *DeliveryRepo) DeleteDelivery(uid string) error {
	sql := `DELETE FROM deliveries WHERE order_id = $1`

	_, err := r.pg.Pool.Exec(context.Background(), sql, uid)
	if err != nil {
		log.Errorf("Failed to delete delivery from db: %v", err)
		return err
	}
	log.Infoln("Delivery deleted from db")
	return nil
}

func (r *DeliveryRepo) DeleteAllDeliveries() error {
	sql := `DELETE FROM deliveries`

	_, err := r.pg.Pool.Exec(context.Background(), sql)
	if err != nil {
		log.Errorf("Failed to delete all deliveries from db: %v", err)
		return err
	}
	log.Infoln("All deliveries deleted from db")
	return nil
}
