package memory

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
)

type ItemRepo struct {
	Items []model.Item
	sm    sync.RWMutex
}

func NewItemRepo(items []model.Item) *ItemRepo {
	return &ItemRepo{
		Items: items,
	}
}

func (it *ItemRepo) CreateItem(item model.Item) error {
	it.sm.Lock()
	defer it.sm.Unlock()

	it.Items = append(it.Items, item)
	return nil
}

func (it *ItemRepo) CreateItems(items []model.Item) error {
	it.sm.Lock()
	defer it.sm.Unlock()

	it.Items = append(it.Items, items...)
	return nil
}

func (it *ItemRepo) GetItemsByID(id string) ([]model.Item, error) {
	it.sm.RLock()
	defer it.sm.RUnlock()

	items := make([]model.Item, 0)
	for _, item := range it.Items {
		if item.OrderID == id {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("items not found in memory")
	} else {
		return items, nil
	}
}

func (it *ItemRepo) GetAllItems() []model.Item {
	return it.Items
}

func (it *ItemRepo) DeleteItems(uid string) error {
	it.sm.Lock()
	defer it.sm.Unlock()

	var found bool

	for i, item := range it.Items {
		if item.OrderID == uid {
			it.Items = append(it.Items[:i], it.Items[i+1:]...)
			found = true
		}
	}
	if !found {
		return fmt.Errorf("items not found in memory")
	} else {
		return nil
	}
}

func (it *ItemRepo) DeleteAllItems() error {
	it.sm.Lock()
	defer it.sm.Unlock()

	it.Items = nil
	log.Infoln("All items deleted from memory")
	return nil
}
