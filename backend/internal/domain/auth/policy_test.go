package auth_test

import (
	"testing"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
)

func TestValidatePassword(t *testing.T) {
	if err := auth.ValidatePassword("12345678"); err != nil {
		t.Fatalf("expected 8-char password to be valid, got %v", err)
	}

	err := auth.ValidatePassword("short")
	if err == nil {
		t.Fatal("expected error for password shorter than 8 chars")
	}
	var domainErr *domainerrors.DomainError
	if de, ok := err.(*domainerrors.DomainError); ok {
		domainErr = de
	} else {
		t.Fatalf("expected *DomainError, got %T", err)
	}
	if domainErr.Code != domainerrors.CodeValidation {
		t.Fatalf("expected CodeValidation, got %q", domainErr.Code)
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := auth.HashPassword("mi-password-segura")
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}
	if hash == "" || hash == "mi-password-segura" {
		t.Fatalf("expected a bcrypt hash, got %q", hash)
	}

	if !auth.CheckPassword(hash, "mi-password-segura") {
		t.Fatal("expected CheckPassword to succeed with the correct password")
	}
	if auth.CheckPassword(hash, "otra-password") {
		t.Fatal("expected CheckPassword to fail with the wrong password")
	}
}

func TestRefreshTokenTTL(t *testing.T) {
	adminTTL := 720 * time.Hour
	vendedoraTTL := 8 * time.Hour

	if got := auth.RefreshTokenTTL("admin", adminTTL, vendedoraTTL); got != adminTTL {
		t.Errorf("expected admin TTL %v, got %v", adminTTL, got)
	}
	if got := auth.RefreshTokenTTL("vendedora", adminTTL, vendedoraTTL); got != vendedoraTTL {
		t.Errorf("expected vendedora TTL %v, got %v", vendedoraTTL, got)
	}
}
