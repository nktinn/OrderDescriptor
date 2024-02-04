package sub

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/model"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo"
	log "github.com/sirupsen/logrus"
)

type NatsSubscriber struct {
	memrepo *repo.MemoryRepo
	sc      stan.Conn
	sub     map[string]stan.Subscription
}

func NewSubscriber(sc stan.Conn, memrepo *repo.MemoryRepo) *NatsSubscriber {
	return &NatsSubscriber{
		memrepo: memrepo,
		sc:      sc,
		sub:     make(map[string]stan.Subscription),
	}
}

func (n *NatsSubscriber) Subscribe(subject string) error {
	sub, err := n.sc.Subscribe(subject, func(msg *stan.Msg) {
		log.Infoln("Received message from nats-streaming")
		order, err := n.Validate(msg.Data)
		if err != nil {
			log.Errorf("Validation error: %v", err)
			log.Errorf("Subscribe method finished with error")
			return
		}
		log.Infoln("Order validated")
		err = n.memrepo.CreateOrder(order)
		if err != nil {
			log.Errorf("Subscribe method finished with error")
			return
		}
		log.Infoln("Order created. UID:", order.OrderUID)
		log.Infoln("Subscribe method finished successfully")
	})
	if err != nil {
		return err
	}
	n.sub[subject] = sub
	return nil
}

func (n *NatsSubscriber) Validate(data []byte) (model.Order, error) {
	var order model.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (n *NatsSubscriber) Unsubscribe(subject string) error {
	if sub, ok := n.sub[subject]; ok {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
		delete(n.sub, subject)
		return nil
	}
	return fmt.Errorf("subscriber not found")
}

func (n *NatsSubscriber) UnsubscribeAll() error {
	for _, sub := range n.sub {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
	}
	return nil
}
