package auth

import "time"

// RefreshToken is the domain representation of an opaque refresh token
// (the plaintext value is never persisted, only its SHA-256 hash).
type RefreshToken struct {
	ID        int64
	UsuarioID int64
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
}

// Expired reports whether the token is past its expiry at now.
func (t RefreshToken) Expired(now time.Time) bool {
	return now.After(t.ExpiresAt)
}

// Revoked reports whether the token has already been revoked.
func (t RefreshToken) Revoked() bool {
	return t.RevokedAt != nil
}

// Session is the pair of tokens returned to the client after a successful
// login or refresh.
type Session struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string // plaintext; only ever held in memory before hashing/cookie-setting
	RefreshExpiresAt time.Time
}
