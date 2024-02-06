package repo

import (
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo/memory"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo/pgdb"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
)

type Payment interface {
	CreatePayment(payment model.Payment) error
	GetPayment(id string) (*model.Payment, error)
	GetAllPayments() []model.Payment
	DeletePayment(uid string) error
	DeleteAllPayments() error
}

type Delivery interface {
	CreateDelivery(delivery model.Delivery) error
	GetDelivery(id string) (*model.Delivery, error)
	GetAllDeliveries() []model.Delivery
	DeleteDelivery(uid string) error
	DeleteAllDeliveries() error
}

type Order interface {
	CreateOrder(order model.Order) error
	GetOrder(uid string) (*model.Order, error)
	GetAllOrders() []model.Order
	DeleteOrder(uid string) error
	DeleteAllOrders() error
}

type Item interface {
	CreateItem(item model.Item) error
	CreateItems(items []model.Item) error
	GetItemsByID(id string) ([]model.Item, error)
	GetAllItems() []model.Item
	DeleteItems(uid string) error
	DeleteAllItems() error
}

type Repositories struct {
	Payment
	Delivery
	Order
	Item
}

// Database repositories constructor
func NewDBRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Payment:  pgdb.NewPaymentRepo(pg),
		Delivery: pgdb.NewDeliveryRepo(pg),
		Order:    pgdb.NewOrderRepo(pg),
		Item:     pgdb.NewItemRepo(pg),
	}
}

// In-Memory repositories constructor
func NewMemoryRepositories(orders []model.Order, deliveries []model.Delivery, payments []model.Payment, items []model.Item) *Repositories {
	return &Repositories{
		Payment:  memory.NewPaymentRepo(payments),
		Delivery: memory.NewDeliveryRepo(deliveries),
		Order:    memory.NewOrderRepo(orders),
		Item:     memory.NewItemRepo(items),
	}
}
