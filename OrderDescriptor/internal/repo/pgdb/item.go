package pgdb

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
)

type ItemRepo struct {
	pg *postgres.Postgres
}

func NewItemRepo(pg *postgres.Postgres) *ItemRepo {
	return &ItemRepo{pg}
}

func (r *ItemRepo) CreateItem(item model.Item) error {
	sql := `
  	INSERT INTO items (order_id, chrt_id, track_number, price, rid, item_name, sale, item_size, total_price, nm_id, brand, status)
  	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
 	`

	_, err := r.pg.Pool.Exec(context.Background(), sql,
		item.OrderID,
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
		return fmt.Errorf("failed to insert new item: %v", err)
	}
	return nil
}

func (r *ItemRepo) CreateItems(items []model.Item) error {
	for _, item := range items {
		err := r.CreateItem(item)
		if err != nil {
			return fmt.Errorf("failed to insert new items: %v", err)
		}
	}
	return nil
}

func (r *ItemRepo) GetItemsByID(id string) ([]model.Item, error) {
	items := make([]model.Item, 0)

	query := `SELECT * FROM items WHERE order_id = $1`
	rows, err := r.pg.Pool.Query(context.Background(), query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		item := model.Item{}
		err = rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan items: %v", err)
		}
		items = append(items, item)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	return items, nil
}

func (r *ItemRepo) GetAllItems() []model.Item {
	items := make([]model.Item, 0)

	query := `SELECT * FROM items`
	rows, err := r.pg.Pool.Query(context.Background(), query)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		item := model.Item{}
		err = rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			log.Errorf("Failed to write item to memory: %v", err)
			return nil
		}
		items = append(items, item)
	}
	if err != nil {
		log.Errorf("Failed to get items: %v", err)
		return nil
	}

	log.Infoln("Items wrote to memory")
	return items
}

func (r *ItemRepo) DeleteItems(uid string) error {
	sql := `DELETE FROM items WHERE order_id = $1`
	_, err := r.pg.Pool.Exec(context.Background(), sql, uid)
	if err != nil {
		log.Errorf("Failed to delete items from db: %v", err)
		return err
	}
	log.Infoln("Items deleted from db")
	return nil
}

func (r *ItemRepo) DeleteAllItems() error {
	sql := `DELETE FROM items`

	_, err := r.pg.Pool.Exec(context.Background(), sql)
	if err != nil {
		log.Errorf("Failed to delete all items from db: %v", err)
		return err
	}
	log.Infoln("All items deleted from db")
	return nil
}
