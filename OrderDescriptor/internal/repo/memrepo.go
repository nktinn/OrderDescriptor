package repo

import (
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

type MemoryRepo struct {
	Payments   []model.Payment
	Deliveries []model.Delivery
	Items      []model.Item
	Orders     []model.Order

	repo *Repositories
	sm   sync.RWMutex
}

func NewMemoryRepo(repo *Repositories) *MemoryRepo {
	return &MemoryRepo{
		Payments:   repo.DBPayment.GetAllPayments(),
		Deliveries: repo.DBDelivery.GetAllDeliveries(),
		Items:      repo.DBItem.GetAllItems(),
		Orders:     repo.DBOrder.GetAllOrders(),
		repo:       repo,
	}
}

// Payments

func (mem *MemoryRepo) CreatePayment(payment model.Payment) error {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	err := mem.repo.CreatePayment(payment)
	if err != nil {
		log.Errorf("Failed to insert payment into db: %v", err)
		return err
	}
	log.Infoln("Payment added to db")

	mem.Payments = append(mem.Payments, payment)
	log.Infoln("Payment added to memory")
	return nil
}

func (mem *MemoryRepo) GetPayment(id string) *model.Payment {
	mem.sm.RLock()
	defer mem.sm.RUnlock()
	for _, payment := range mem.Payments {
		if payment.OrderID == id {
			log.Infoln("Payment found in memory")
			return &payment
		}
	}
	log.Infoln("Payment not found in memory")
	return nil
}

func (mem *MemoryRepo) DeletePayment(uid string) {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	for i, payment := range mem.Payments {
		if payment.OrderID == uid {
			mem.Payments = append(mem.Payments[:i], mem.Payments[i+1:]...)
			break
		}
	}
}

// Deliveries

func (mem *MemoryRepo) CreateDelivery(delivery model.Delivery) error {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	err := mem.repo.CreateDelivery(delivery)
	if err != nil {
		log.Errorf("Failed to insert delivery into db: %v", err)
		return err
	}
	log.Infoln("Delivery added to db")

	mem.Deliveries = append(mem.Deliveries, delivery)
	log.Infoln("Delivery added to memory")

	return nil
}

func (mem *MemoryRepo) GetDelivery(id string) *model.Delivery {
	mem.sm.RLock()
	defer mem.sm.RUnlock()

	for _, delivery := range mem.Deliveries {
		if delivery.OrderID == id {
			log.Infoln("Delivery found in memory")
			return &delivery
		}
	}
	log.Infoln("Delivery not found in memory")
	return nil
}

func (mem *MemoryRepo) DeleteDelivery(uid string) {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	for i, delivery := range mem.Deliveries {
		if delivery.OrderID == uid {
			mem.Deliveries = append(mem.Deliveries[:i], mem.Deliveries[i+1:]...)
			break
		}
	}
}

// Items

func (mem *MemoryRepo) CreateItems(items []model.Item) error {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	for _, item := range items {
		err := mem.repo.CreateItem(item)
		if err != nil {
			log.Errorf("Failed to insert item into db: %v", err)
			return err
		}
		mem.Items = append(mem.Items, item)
	}
	log.Infoln("Items added to memory")
	log.Infoln("Items added to db")
	return nil
}

func (mem *MemoryRepo) GetItemsByID(id string) []model.Item {
	mem.sm.RLock()
	defer mem.sm.RUnlock()

	items := make([]model.Item, 0)
	for _, item := range mem.Items {
		if item.OrderID == id {
			items = append(items, item)
		}
	}
	return items
}

func (mem *MemoryRepo) DeleteItems(uid string) {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	for i, item := range mem.Items {
		if item.OrderID == uid {
			mem.Items = append(mem.Items[:i], mem.Items[i+1:]...)
		}
	}
}

// Orders

func (mem *MemoryRepo) CreateOrder(order model.Order) error {
	mem.sm.Lock()

	err := mem.repo.CreateOrder(order)
	if err != nil {
		log.Errorf("Failed to create order: %v", err)
		return err
	}
	log.Infoln("Order added to db")

	// Add to memory
	mem.Orders = append(mem.Orders, order)
	log.Infoln("Order added to memory")

	mem.sm.Unlock()

	err = mem.CreateDelivery(*order.Delivery)
	if err != nil {
		log.Errorf("Failed to create delivery: %v", err)
		log.Infoln("Returning applied changes...")
		mem.DeleteOrder(order.OrderUID)
		return err
	}
	err = mem.CreatePayment(*order.Payment)
	if err != nil {
		log.Errorf("Failed to create payment: %v", err)
		log.Infoln("Returning applied changes...")
		mem.DeleteOrder(order.OrderUID)
		return err
	}
	err = mem.CreateItems(order.Items)
	if err != nil {
		log.Errorf("Failed to create items: %v", err)
		log.Infoln("Returning applied changes...")
		mem.DeleteOrder(order.OrderUID)
		return err
	}

	log.Infoln("Order created")
	return nil
}

func (mem *MemoryRepo) GetOrder(uid string) *model.Order {
	mem.sm.RLock()
	defer mem.sm.RUnlock()

	for _, order := range mem.Orders {
		if order.OrderUID == uid {
			log.Infoln("Order found in memory")
			return &order
		}
	}
	log.Errorf("Order not found in memory")
	return nil
}

func (mem *MemoryRepo) GetFullOrder(uid string) *model.Order {
	mem.sm.RLock()
	defer mem.sm.RUnlock()

	order := mem.GetOrder(uid)
	if order == nil {
		log.Errorf("Order not found while getting full order")
		return nil
	}
	delivery := mem.GetDelivery(order.OrderUID)
	if delivery == nil {
		log.Errorf("Delivery not found while getting full order")
		return nil
	}
	payment := mem.GetPayment(order.OrderUID)
	if payment == nil {
		log.Errorf("Payment not found while getting full order")
		return nil
	}
	items := mem.GetItemsByID(order.OrderUID)
	if items == nil {
		log.Errorf("Items not found while getting full order")
		return nil
	}

	order.Payment = payment
	order.Delivery = delivery
	order.Items = items

	return order
}

func (mem *MemoryRepo) DeleteOrder(uid string) {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	mem.repo.DeleteOrder(uid)

	for i, order := range mem.Orders {
		if order.OrderUID == uid {
			mem.Orders = append(mem.Orders[:i], mem.Orders[i+1:]...)
			break
		}
	}

	mem.DeleteDelivery(uid)
	mem.DeletePayment(uid)
	mem.DeleteItems(uid)
}

// DeleteAll

func (mem *MemoryRepo) DeleteAll() error {
	mem.sm.Lock()
	defer mem.sm.Unlock()

	var err error
	for _, order := range mem.Orders {
		err = mem.repo.DeleteOrder(order.OrderUID)
		if err != nil {
			return err
		}
	}

	mem.Orders = nil
	mem.Payments = nil
	mem.Deliveries = nil
	mem.Items = nil

	return nil
}
