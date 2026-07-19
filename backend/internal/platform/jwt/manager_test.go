package jwt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	m := jwt.New("test-secret", 15*time.Minute)

	token, expiresAt, err := m.GenerateAccessToken(42, "admin", 1)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if !expiresAt.After(time.Now()) {
		t.Fatalf("expected expiresAt in the future, got %v", expiresAt)
	}

	claims, err := m.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("unexpected error parsing token: %v", err)
	}
	if claims.UserID != 42 || claims.Rol != "admin" || claims.SedeID != 1 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseExpiredToken(t *testing.T) {
	m := jwt.New("test-secret", -1*time.Minute)

	token, _, err := m.GenerateAccessToken(1, "vendedora", 1)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = m.ParseAccessToken(token)
	if !errors.Is(err, jwt.ErrExpiredToken) {
		t.Fatalf("expected ErrExpiredToken, got %v", err)
	}
}

func TestParseTokenWithWrongSecret(t *testing.T) {
	issuer := jwt.New("secret-a", 15*time.Minute)
	token, _, err := issuer.GenerateAccessToken(1, "admin", 1)
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	verifier := jwt.New("secret-b", 15*time.Minute)
	_, err = verifier.ParseAccessToken(token)
	if !errors.Is(err, jwt.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}

func TestParseGarbageToken(t *testing.T) {
	m := jwt.New("test-secret", 15*time.Minute)
	if _, err := m.ParseAccessToken("not-a-jwt"); !errors.Is(err, jwt.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}
}
