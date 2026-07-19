package db_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/db"
)

// TestConnectAndHealthCheck is a lightweight integration test that runs
// against the local dev Postgres (DATABASE_URL, e.g. via `make db-up`). It
// skips instead of failing when no database is reachable, so `make test`
// stays green without Docker running.
func TestConnectAndHealthCheck(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping db integration test")
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, dbURL, logger)
	if err != nil {
		t.Skipf("could not connect to %s, skipping: %v", dbURL, err)
	}
	defer pool.Close()

	if err := db.HealthCheck(ctx, pool); err != nil {
		t.Fatalf("HealthCheck failed on a pool that just connected: %v", err)
	}
}
