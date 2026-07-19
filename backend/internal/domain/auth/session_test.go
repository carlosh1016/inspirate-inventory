package auth_test

import (
	"testing"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
)

func TestRefreshTokenExpiredAndRevoked(t *testing.T) {
	now := time.Now()

	notExpired := auth.RefreshToken{ExpiresAt: now.Add(time.Hour)}
	if notExpired.Expired(now) {
		t.Error("expected token to not be expired")
	}

	expired := auth.RefreshToken{ExpiresAt: now.Add(-time.Hour)}
	if !expired.Expired(now) {
		t.Error("expected token to be expired")
	}

	notRevoked := auth.RefreshToken{}
	if notRevoked.Revoked() {
		t.Error("expected token to not be revoked")
	}

	revokedAt := now
	revoked := auth.RefreshToken{RevokedAt: &revokedAt}
	if !revoked.Revoked() {
		t.Error("expected token to be revoked")
	}
}

func TestPasswordResetExpiredAndUsed(t *testing.T) {
	now := time.Now()

	valid := auth.PasswordReset{ExpiresAt: now.Add(time.Hour)}
	if valid.Expired(now) {
		t.Error("expected reset to not be expired")
	}
	if valid.Used() {
		t.Error("expected reset to not be used")
	}

	usedAt := now
	used := auth.PasswordReset{ExpiresAt: now.Add(time.Hour), UsedAt: &usedAt}
	if !used.Used() {
		t.Error("expected reset to be used")
	}
}
