package auth

import "time"

// PasswordReset is the domain representation of a single-use password reset
// token (the plaintext value is never persisted, only its SHA-256 hash).
type PasswordReset struct {
	ID        int64
	UsuarioID int64
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
}

// Expired reports whether the token is past its expiry at now.
func (p PasswordReset) Expired(now time.Time) bool {
	return now.After(p.ExpiresAt)
}

// Used reports whether the token has already been consumed.
func (p PasswordReset) Used() bool {
	return p.UsedAt != nil
}
