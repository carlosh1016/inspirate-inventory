// Package idempotencykeys is the persistence port for idempotency_keys.
// Rows are write-once (except the periodic expiry cleanup, which deletes
// rather than updates).
package idempotencykeys

import (
	"context"
	"errors"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no idempotency key matches the lookup (or it
// exists but already expired — GetByKey filters expires_at > NOW() at the
// SQL level, so an expired row behaves identically to a missing one).
var ErrNotFound = errors.New("idempotency key not found")

// InsertParams is the full set of columns for a new idempotency key row.
type InsertParams struct {
	Key            string
	UsuarioID      int64
	Endpoint       string
	ResponseStatus int32
	ResponseBody   []byte
	RequestHash    string
	ExpiresAt      time.Time
}

// Repository is the persistence port for idempotency_keys, consumed by
// usecase/ventas. NewPostgres accepts generated.DBTX, so Insert can run
// inside CreateVenta's own transaction, alongside the venta it caches the
// response for.
type Repository interface {
	GetByKey(ctx context.Context, key string) (generated.IdempotencyKey, error)
	Insert(ctx context.Context, params InsertParams) error
	DeleteExpired(ctx context.Context) error
}
