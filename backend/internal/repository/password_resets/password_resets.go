// Package passwordresets is the persistence port for the password_resets
// table (opaque single-use reset tokens).
package passwordresets

import (
	"context"
	"errors"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no password reset matches the lookup.
var ErrNotFound = errors.New("password reset not found")

// Repository is the persistence port for password_resets, consumed by usecases.
type Repository interface {
	Insert(ctx context.Context, usuarioID int64, tokenHash string, expiresAt time.Time) (generated.PasswordReset, error)
	GetByHash(ctx context.Context, tokenHash string) (generated.PasswordReset, error)
	MarkUsed(ctx context.Context, id int64) error
}
