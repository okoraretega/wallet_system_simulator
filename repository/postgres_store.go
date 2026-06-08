package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.New("Failed to parse connection string")
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.New("Unable to connect to database")
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, errors.New("Failed to ping database")
	}

	return &PostgresStore{db: pool}, nil
}

func (s *PostgresStore) Close() {
	s.db.Close()
}
