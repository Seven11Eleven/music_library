package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Seven11Eleven/music_library/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

type Postgres interface {
	DB() *pgxpool.Pool
	CloseDB() error
}

type postgres struct {
	pool *pgxpool.Pool
}

func NewDB(cfg config.Config) (*postgres, error) {
	log.Info("Initializing database connection...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Infof("Connecting to database at: %s", dbURL)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
		return nil, err
	}

	config.MaxConns = 10
	log.Infof("Setting max connections to: %d", config.MaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil, err
	}

	log.Info("Database connection successfully established")
	return &postgres{pool: pool}, nil
}

func (p *postgres) DB() *pgxpool.Pool {
	return p.pool
}

func (p *postgres) CloseDB() error {
	if p.pool == nil {
		log.Warn("CloseDB called but the pool is nil")
		return errors.New("not connected")
	}

	log.Info("Closing connection pool to the database")
	p.pool.Close()
	log.Info("Database connection pool closed successfully")
	return nil
}
