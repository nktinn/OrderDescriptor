package app

import (
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/config"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/controller"
	natsSub "github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/controller/nats"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/internal/repo"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/pkg/postgres"
	"github.com/nktinn/OrderDescriptor/pkg/nats"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	// Logger setup
	SetLogger()

	log.Infoln("Starting application...")
	// Configuration
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Unable to read config: %v", err)
	}
	log.Infoln("Read config")

	// Postgres connection
	pg, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Infoln("Connected to the database")
	defer pg.Close()

	// New repositories
	repositories := repo.NewRepositories(pg)
	memrepos := repo.NewMemoryRepo(repositories)
	log.Infoln("Created repositories")

	// Start server
	go func() {
		controller.StartServer(cfg.HTTP, repositories, memrepos)
	}()
	log.Infoln("Started server")

	// Connect to nats-streaming
	natsConn := nats.NewConnection(cfg.Nats)
	if natsConn == nil {
		log.Fatalf("Unable to connect to nats-streaming")
	} else {
		log.Infoln("Connected to nats-streaming")
	}
	defer log.Infoln("Disconnected from nats-streaming")
	defer natsConn.CloseConnection()

	// Subscribe to nats-streaming
	natsSubscriber := natsSub.NewSubscriber(natsConn.SC, memrepos)
	if natsSubscriber == nil {
		log.Error("Unable to connect to nats-streaming")
	}
	log.Infoln("Created subscriber")
	if err = natsSubscriber.Subscribe(cfg.Nats.Subject); err != nil {
		log.Errorf("Error while subscribing to %s: %v", cfg.Nats.Subject, err)
	} else {
		log.Infof("Subscribed to %s", cfg.Nats.Subject)

		defer func() {
			if err = natsSubscriber.Unsubscribe(cfg.Nats.Subject); err != nil {
				log.Errorf("Error while unsubscribing from %s: %v", cfg.Nats.Subject, err)
			} else {
				log.Infof("Unsubscribed from %s", cfg.Nats.Subject)
			}
		}()
	}

	// Ping database connection
	go pg.Ping()

	// Shutdown with Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	defer log.Infoln("Shutting down application...")
}
