package config_test

import (
	"strings"
	"testing"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/config"
)

func setBaseEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	t.Setenv("JWT_SECRET", "dev-secret")
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("LOG_LEVEL", "info")
}

func TestLoadSuccess(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("PORT", "9090")
	t.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 9090 {
		t.Errorf("expected Port=9090, got %d", cfg.Port)
	}
	if len(cfg.CORSAllowedOrigins) != 2 {
		t.Errorf("expected 2 CORS origins, got %v", cfg.CORSAllowedOrigins)
	}
	if cfg.JWTAccessTTLAdmin.String() != "24h0m0s" {
		t.Errorf("expected default JWTAccessTTLAdmin=24h, got %s", cfg.JWTAccessTTLAdmin)
	}
	if cfg.JWTAccessTTLVendedora.String() != "10m0s" {
		t.Errorf("expected default JWTAccessTTLVendedora=10m, got %s", cfg.JWTAccessTTLVendedora)
	}
}

func TestLoadMissingRequiredField(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("DATABASE_URL", "")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for missing DATABASE_URL, got nil")
	}
}

func TestLoadInvalidEnvironment(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("ENVIRONMENT", "staging")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for invalid ENVIRONMENT, got nil")
	}
	if !strings.Contains(err.Error(), "ENVIRONMENT") {
		t.Errorf("expected error to mention ENVIRONMENT, got %q", err.Error())
	}
}

func TestLoadInvalidLogLevel(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("LOG_LEVEL", "verbose")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for invalid LOG_LEVEL, got nil")
	}
}

func TestLoadProductionRequiresLongJWTSecret(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", "too-short")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for short JWT_SECRET in production, got nil")
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Errorf("expected error to mention JWT_SECRET, got %q", err.Error())
	}
}

func TestLoadProductionAcceptsLongJWTSecret(t *testing.T) {
	setBaseEnv(t)
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))

	if _, err := config.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
