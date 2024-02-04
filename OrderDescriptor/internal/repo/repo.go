package repo

import (
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo/pgdb"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
)

type DBPayment interface {
	CreatePayment(payment model.Payment) error
	GetPayment(id string) (*model.Payment, error)
	GetAllPayments() []model.Payment
	DeletePayment(uid string) error
}

type DBDelivery interface {
	CreateDelivery(delivery model.Delivery) error
	GetDelivery(id string) (*model.Delivery, error)
	GetAllDeliveries() []model.Delivery
	DeleteDelivery(uid string) error
}

type DBOrder interface {
	CreateOrder(order model.Order) error
	GetOrder(uid string) (*model.Order, error)
	GetAllOrders() []model.Order
	DeleteOrder(uid string) error
}

type DBItem interface {
	CreateItem(item model.Item) error
	CreateItems(items []model.Item) error
	GetItemsByID(id string) ([]model.Item, error)
	GetAllItems() []model.Item
	DeleteItems(uid string) error
}

type Repositories struct {
	DBPayment
	DBDelivery
	DBOrder
	DBItem
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	deliveryRepo := pgdb.NewDeliveryRepo(pg)
	paymentRepo := pgdb.NewPaymentRepo(pg)
	itemRepo := pgdb.NewItemRepo(pg)
	orderRepo := pgdb.NewOrderRepo(pg, deliveryRepo, paymentRepo, itemRepo)

	return &Repositories{
		DBPayment:  paymentRepo,
		DBDelivery: deliveryRepo,
		DBOrder:    orderRepo,
		DBItem:     itemRepo,
	}
}
