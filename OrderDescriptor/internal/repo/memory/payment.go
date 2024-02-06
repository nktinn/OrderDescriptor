package memory

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
)

type PaymentRepo struct {
	Payments []model.Payment
	sm       sync.RWMutex
}

func NewPaymentRepo(payments []model.Payment) *PaymentRepo {
	return &PaymentRepo{
		Payments: payments,
	}
}

func (pay *PaymentRepo) CreatePayment(payment model.Payment) error {
	pay.sm.Lock()
	defer pay.sm.Unlock()

	pay.Payments = append(pay.Payments, payment)
	log.Infoln("Payment added to memory")
	return nil
}

func (pay *PaymentRepo) GetPayment(id string) (*model.Payment, error) {
	pay.sm.RLock()
	defer pay.sm.RUnlock()

	for _, payment := range pay.Payments {
		if payment.OrderID == id {
			log.Infoln("Payment found in memory")
			return &payment, nil
		}
	}
	log.Infoln("Payment not found in memory")
	return nil, fmt.Errorf("payment not found in memory")
}

func (pay *PaymentRepo) GetAllPayments() []model.Payment {
	return pay.Payments
}

func (pay *PaymentRepo) DeletePayment(uid string) error {
	pay.sm.Lock()
	defer pay.sm.Unlock()

	for i, payment := range pay.Payments {
		if payment.OrderID == uid {
			pay.Payments = append(pay.Payments[:i], pay.Payments[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("payment not found in memory")
}

func (pay *PaymentRepo) DeleteAllPayments() error {
	pay.sm.Lock()
	defer pay.sm.Unlock()

	pay.Payments = nil
	log.Infoln("All payments deleted from memory")
	return nil
}
