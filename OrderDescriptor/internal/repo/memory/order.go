package memory

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
)

type OrderRepo struct {
	Orders []model.Order
	sm     sync.RWMutex
}

func NewOrderRepo(orders []model.Order) *OrderRepo {
	return &OrderRepo{
		Orders: orders,
	}
}

func (ord *OrderRepo) CreateOrder(order model.Order) error {
	ord.sm.Lock()
	defer ord.sm.Unlock()

	ord.Orders = append(ord.Orders, order)
	return nil
}

func (ord *OrderRepo) GetOrder(uid string) (*model.Order, error) {
	ord.sm.RLock()
	defer ord.sm.RUnlock()

	for _, order := range ord.Orders {
		if order.OrderUID == uid {
			log.Infoln("Order found in memory")
			return &order, nil
		}
	}
	log.Errorf("Order not found in memory")
	return nil, fmt.Errorf("order not found in memory")
}

func (ord *OrderRepo) GetAllOrders() []model.Order {
	return ord.Orders
}

func (ord *OrderRepo) DeleteOrder(uid string) error {
	ord.sm.Lock()
	defer ord.sm.Unlock()

	for i, order := range ord.Orders {
		if order.OrderUID == uid {
			ord.Orders = append(ord.Orders[:i], ord.Orders[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("order not found in memory")
}

func (ord *OrderRepo) DeleteAllOrders() error {
	ord.sm.Lock()
	defer ord.sm.Unlock()

	ord.Orders = nil
	log.Infoln("All orders deleted from memory")
	return nil
}

/*
func (ord *OrderRepo) GetFullOrder(uid string) *model.Order {
	ord.sm.RLock()
	defer ord.sm.RUnlock()

	order := ord.GetOrder(uid)
	if order == nil {
		log.Errorf("Order not found while getting full order")
		return nil
	}
	delivery := ord.GetDelivery(order.OrderUID)
	if delivery == nil {
		log.Errorf("Delivery not found while getting full order")
		return nil
	}
	payment := ord.GetPayment(order.OrderUID)
	if payment == nil {
		log.Errorf("Payment not found while getting full order")
		return nil
	}
	items := ord.GetItemsByID(order.OrderUID)
	if items == nil {
		log.Errorf("Items not found while getting full order")
		return nil
	}

	order.Payment = payment
	order.Delivery = delivery
	order.Items = items

	return order
}*/
