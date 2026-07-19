// Package db manages the PostgreSQL connection pool (pgx).
package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxConns        = 20
	minConns        = 2
	maxConnLifetime = 30 * time.Minute
	maxConnIdleTime = 5 * time.Minute
)

// Connect builds and validates a pgx connection pool for dbURL.
func Connect(ctx context.Context, dbURL string, logger *slog.Logger) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database url: %w", err)
	}

	cfg.MaxConns = maxConns
	cfg.MinConns = minConns
	cfg.MaxConnLifetime = maxConnLifetime
	cfg.MaxConnIdleTime = maxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	logger.Info("connected to database", "max_conns", maxConns, "min_conns", minConns)
	return pool, nil
}

// HealthCheck pings the pool; used by GET /health.
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}
