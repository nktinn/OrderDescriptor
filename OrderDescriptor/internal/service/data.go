package service

import (
	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo"
)

type OrderService struct {
	dbRepos  *repo.Repositories
	memRepos *repo.Repositories
}

func NewOrderService(dbRepos *repo.Repositories, memRepos *repo.Repositories) *OrderService {
	return &OrderService{
		dbRepos:  dbRepos,
		memRepos: memRepos,
	}
}

func (ord *OrderService) InitMemory() {
	orders := ord.dbRepos.GetAllOrders()
	deliveries := ord.dbRepos.GetAllDeliveries()
	payments := ord.dbRepos.GetAllPayments()
	items := ord.dbRepos.GetAllItems()

	ord.memRepos = repo.NewMemoryRepositories(orders, deliveries, payments, items)
}

func (ord *OrderService) CreateOrder(order model.Order) error {
	err := ord.dbRepos.Order.CreateOrder(order)
	if err != nil {
		log.Errorf("Failed to create order: %v", err)
		return err
	}
	log.Infoln("Order added to db")
	ord.memRepos.Order.CreateOrder(order)
	log.Infoln("Order added to memory")

	err = ord.dbRepos.CreateDelivery(*order.Delivery)
	if err != nil {
		log.Errorf("Failed to create delivery: %v", err)
		log.Infoln("Returning applied changes...")
		errOrder := ord.dbRepos.DeleteOrder(order.OrderUID)
		if errOrder != nil {
			log.Errorf("Failed to delete order: %v", errOrder)
		}
		ord.memRepos.DeleteOrder(order.OrderUID)
		return err
	}
	log.Infoln("Delivery added to db")
	ord.memRepos.CreateDelivery(*order.Delivery)
	log.Infoln("Delivery added to memory")

	err = ord.dbRepos.CreatePayment(*order.Payment)
	if err != nil {
		log.Errorf("Failed to create payment: %v", err)
		log.Infoln("Returning applied changes...")
		errOrder := ord.dbRepos.DeleteOrder(order.OrderUID)
		if errOrder != nil {
			log.Errorf("Failed to delete order: %v", errOrder)
		}
		ord.memRepos.DeleteOrder(order.OrderUID)
		ord.memRepos.DeleteDelivery(order.OrderUID)
		return err
	}
	log.Infoln("Payment added to db")
	ord.memRepos.CreatePayment(*order.Payment)
	log.Infoln("Payment added to memory")

	err = ord.dbRepos.CreateItems(order.Items)
	if err != nil {
		log.Errorf("Failed to create items: %v", err)
		log.Infoln("Returning applied changes...")
		errOrder := ord.dbRepos.DeleteOrder(order.OrderUID)
		if errOrder != nil {
			log.Errorf("Failed to delete order: %v", errOrder)
		}
		ord.memRepos.DeleteOrder(order.OrderUID)
		ord.memRepos.DeleteDelivery(order.OrderUID)
		ord.memRepos.DeletePayment(order.OrderUID)
		return err
	}
	log.Infoln("Items added to db")
	ord.memRepos.CreateItems(order.Items)
	log.Infoln("Items added to memory")

	log.Infoln("Order created")
	return nil
}

func (ord *OrderService) GetFullOrder(uid string) (*model.Order, error) {

	order, err := ord.memRepos.GetOrder(uid)
	if err != nil {
		log.Errorf("Failed to get order: %v", err)
		return nil, err
	}
	order.Delivery, _ = ord.memRepos.GetDelivery(uid)
	order.Payment, _ = ord.memRepos.GetPayment(uid)
	order.Items, _ = ord.memRepos.GetItemsByID(uid)

	return order, nil
}

func (ord *OrderService) DeleteAll() error {
	var err error
	for _, order := range ord.memRepos.GetAllOrders() {
		log.Infoln("Deleting order: ", order.OrderUID)
		err = ord.dbRepos.DeleteOrder(order.OrderUID)
		if err != nil {
			log.Errorf("Failed to delete order: %v", err)
			return err
		}
	}

	ord.memRepos.DeleteAllOrders()
	ord.memRepos.DeleteAllDeliveries()
	ord.memRepos.DeleteAllItems()
	ord.memRepos.DeleteAllOrders()

	return nil
}
