package productos_test

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testPool is shared by every test in this package; each test truncates the
// tables it needs at the start (see testenv_test.go) to stay isolated.
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	flag.Parse() // testing.Short() below requires flags to be parsed first.
	ctx := context.Background()

	var container *tcpostgres.PostgresContainer
	if !testing.Short() {
		var pool *pgxpool.Pool
		var err error
		container, pool, err = setupPostgresContainer(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "productos integration tests: could not start postgres testcontainer, skipping: %v\n", err)
		} else {
			testPool = pool
		}
	}

	code := m.Run()

	if testPool != nil {
		testPool.Close()
	}
	if container != nil {
		_ = container.Terminate(ctx)
	}

	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context) (*tcpostgres.PostgresContainer, *pgxpool.Pool, error) {
	container, err := tcpostgres.Run(ctx, "postgres:15-alpine",
		tcpostgres.WithDatabase("inspirate_test"),
		tcpostgres.WithUsername("inspirate"),
		tcpostgres.WithPassword("devpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("starting container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("connection string: %w", err)
	}

	if err := applyMigrations(connStr); err != nil {
		return nil, nil, fmt.Errorf("applying migrations: %w", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("creating pool: %w", err)
	}

	return container, pool, nil
}

// applyMigrations runs db/migrations against the container using goose as a
// library (not the CLI, which pulls in every SQL driver it supports).
func applyMigrations(connStr string) error {
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(sqlDB, "../../../db/migrations")
}

func requireTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testPool == nil {
		t.Skip("postgres testcontainer not available (short mode or Docker unavailable)")
	}
	return testPool
}
