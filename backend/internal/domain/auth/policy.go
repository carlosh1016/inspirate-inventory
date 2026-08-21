// Package auth contains pure domain logic for authentication: password
// policy and session/reset token entities. No I/O, no HTTP, no SQL.
package auth

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

const (
	minPasswordLen = 8
	bcryptCost     = 12
)

// ValidatePassword enforces the password policy: at least 8 characters, no
// other constraints.
func ValidatePassword(password string) error {
	if len(password) < minPasswordLen {
		return domainerrors.NewValidation(
			"Contraseña inválida",
			"La contraseña debe tener al menos 8 caracteres.",
			map[string][]string{"password_nueva": {"Debe tener al menos 8 caracteres."}},
		)
	}
	return nil
}

// HashPassword hashes password with bcrypt (cost 12).
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword reports whether password matches hash.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// RefreshTokenTTL picks the refresh-token lifetime for rol: admins get
// adminTTL (30 days), everyone else gets vendedoraTTL (8 hours).
func RefreshTokenTTL(rol string, adminTTL, vendedoraTTL time.Duration) time.Duration {
	if rol == "admin" {
		return adminTTL
	}
	return vendedoraTTL
}

// AccessTokenTTL picks the access-token lifetime for rol: admins get
// adminTTL (~24h), everyone else gets vendedoraTTL (~10min) — vendedoras
// share computers on the sales floor, so their sessions expire fast.
func AccessTokenTTL(rol string, adminTTL, vendedoraTTL time.Duration) time.Duration {
	if rol == "admin" {
		return adminTTL
	}
	return vendedoraTTL
}
