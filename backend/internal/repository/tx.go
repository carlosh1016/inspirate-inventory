package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithTx runs fn inside a pgx transaction on pool. fn receives the *pgx.Tx
// (which satisfies generated.DBTX, so repositories can be built directly on
// it) so every write inside fn is atomic. Commits on success; rolls back on
// error or panic.
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx) // no-op once the transaction is committed
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
