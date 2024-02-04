package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nktinn/OrderDescriptor/OrderDescriptor/config"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	connString string
	Pool       *pgxpool.Pool
}

func NewPostgresDB(cfg config.PostgresConfig) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	// Postgresql connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pg.connString = connString

	// Connect to the postgres database
	var err error
	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.New(context.Background(), connString)
		if err == nil {
			pg.connAttempts = defaultConnAttempts
			return pg, nil
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}
	pg.connAttempts = defaultConnAttempts
	return nil, err
}

func (p *Postgres) Ping() {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()
	for {
		connErr := p.Pool.Ping(context.Background())
		if connErr != nil {
			log.Errorf("Database connection lost: %v\nTrying to reconnect...", connErr)
			for p.connAttempts > 0 {
				p.Pool, connErr = pgxpool.New(context.Background(), p.connString)
				if connErr == nil {
					p.connAttempts = defaultConnAttempts
					log.Infof("Successfully reconnected to the database")
					break
				}
				log.Printf("Postgres is trying to reconnect, attempts left: %d", p.connAttempts)
				time.Sleep(p.connTimeout)
				p.connAttempts--
			}
			log.Fatalf("Unable to reconnect to database: %v\n", connErr)
		}
		log.Infoln("Connection to the database is still alive")
		<-ticker.C
	}
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
