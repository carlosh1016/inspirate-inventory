// Package refreshtokens is the persistence port for the refresh_tokens
// table (opaque session tokens).
package refreshtokens

import (
	"context"
	"errors"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// ErrNotFound is returned when no refresh token matches the lookup.
var ErrNotFound = errors.New("refresh token not found")

// Repository is the persistence port for refresh_tokens, consumed by usecases.
type Repository interface {
	Insert(ctx context.Context, usuarioID int64, tokenHash, ip, userAgent string, expiresAt time.Time) (generated.RefreshToken, error)
	GetByHash(ctx context.Context, tokenHash string) (generated.RefreshToken, error)
	Revoke(ctx context.Context, id int64) error
	RevokeAllByUser(ctx context.Context, usuarioID int64) error
}
