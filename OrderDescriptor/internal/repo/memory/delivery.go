package memory

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
)

type DeliveryRepo struct {
	Deliveries []model.Delivery
	sm         sync.RWMutex
}

func NewDeliveryRepo(deliveries []model.Delivery) *DeliveryRepo {
	return &DeliveryRepo{
		Deliveries: deliveries,
	}
}

func (del *DeliveryRepo) CreateDelivery(delivery model.Delivery) error {
	del.sm.Lock()
	defer del.sm.Unlock()

	del.Deliveries = append(del.Deliveries, delivery)
	return nil
}

func (del *DeliveryRepo) GetDelivery(id string) (*model.Delivery, error) {
	del.sm.RLock()
	defer del.sm.RUnlock()

	for _, delivery := range del.Deliveries {
		if delivery.OrderID == id {
			log.Infoln("Delivery found in memory")
			return &delivery, nil
		}
	}
	log.Infoln("Delivery not found in memory")
	return nil, fmt.Errorf("delivery not found in memory")
}

func (del *DeliveryRepo) GetAllDeliveries() []model.Delivery {
	return del.Deliveries
}

func (del *DeliveryRepo) DeleteDelivery(uid string) error {
	del.sm.Lock()
	defer del.sm.Unlock()

	for i, delivery := range del.Deliveries {
		if delivery.OrderID == uid {
			del.Deliveries = append(del.Deliveries[:i], del.Deliveries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("delivery not found in memory")
}

func (del *DeliveryRepo) DeleteAllDeliveries() error {
	del.sm.Lock()
	defer del.sm.Unlock()

	del.Deliveries = nil
	log.Infoln("All deliveries deleted from memory")
	return nil
}
