package nats

import (
	"github.com/nats-io/stan.go"

	"github.com/nktinn/OrderDescriptor/OrderDescriptor/config"
)

type Nats struct {
	SC stan.Conn
}

func NewConnection(cfg config.Nats) *Nats {
	// Start nats connection
	SC, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.URL))
	if err != nil {
		return nil
	}
	return &Nats{SC}
}
func (n *Nats) CloseConnection() error {
	err := n.SC.Close()
	if err != nil {
		return err
	}
	return nil
}
