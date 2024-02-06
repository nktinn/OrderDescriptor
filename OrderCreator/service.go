package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nktinn/OrderDescriptor/OrderCreator/publisher"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/config"
	"github.com/nktinn/OrderDescriptor/pkg/nats"
)

const configPath = "OrderCreator/config/config.yaml"

func main() {
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("unable to read config: %v", err)
	}
	log.Infoln("Successfully read config")

	// Connect to the nats-streaming server
	natsConn := nats.NewConnection(cfg.Nats)
	defer natsConn.CloseConnection()

	// Send your JSON data to a specific subject
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			if err = publisher.Publish(natsConn.SC, cfg.Nats.Subject); err != nil {
				log.Infof("Unable to publish to nats-streaming: %v", err)
			} else {
				log.Infof("Successfully published to nats-streaming")
			}
			<-ticker.C
		}
	}()

	// Shutdown with Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Infof("Shutting down...")
}
