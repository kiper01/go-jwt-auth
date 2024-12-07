package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	URL      string
	User     string
	Password string
	PoolSize int
}

func ConnectDB(cfg Config) (*pgxpool.Pool, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s", cfg.User, cfg.Password, cfg.URL)

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.PoolSize)

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}
