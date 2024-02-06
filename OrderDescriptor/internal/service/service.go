package service

import (
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo"
)

type Order interface {
	InitMemory()
	CreateOrder(order model.Order) error
	GetFullOrder(uid string) (*model.Order, error)
	DeleteAll() error
}

type Services struct {
	Order
}

func NewServices(dbRepos *repo.Repositories, memRepos *repo.Repositories) *Services {
	return &Services{
		Order: NewOrderService(dbRepos, memRepos),
	}
}
