// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config holds all application configuration, populated from the environment.
type Config struct {
	Port                   int           `env:"PORT" envDefault:"8080"`
	Environment            string        `env:"ENVIRONMENT" envDefault:"development"`
	LogLevel               string        `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL            string        `env:"DATABASE_URL,required,notEmpty"`
	JWTSecret              string        `env:"JWT_SECRET,required,notEmpty"`
	JWTAccessTTL           time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	JWTRefreshTTLAdmin     time.Duration `env:"JWT_REFRESH_TTL_ADMIN" envDefault:"720h"`
	JWTRefreshTTLVendedora time.Duration `env:"JWT_REFRESH_TTL_VENDEDORA" envDefault:"8h"`
	ResendAPIKey           string        `env:"RESEND_API_KEY"`
	MailFrom               string        `env:"MAIL_FROM" envDefault:"noreply@inspirate.co"`
	FrontendURL            string        `env:"FRONTEND_URL" envDefault:"http://localhost:3000"`
	CORSAllowedOrigins     []string      `env:"CORS_ALLOWED_ORIGINS" envSeparator:","`
}

const minProductionJWTSecretLen = 32

// Load reads a .env file if present (ignoring the error when it isn't —
// expected in production, where variables come from the process
// environment), parses Config from the environment, and validates it.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	switch c.Environment {
	case "development", "production", "test":
	default:
		return fmt.Errorf("invalid ENVIRONMENT %q: must be one of development, production, test", c.Environment)
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid LOG_LEVEL %q: must be one of debug, info, warn, error", c.LogLevel)
	}

	if c.Environment == "production" && len(c.JWTSecret) < minProductionJWTSecretLen {
		return fmt.Errorf("JWT_SECRET must be at least %d characters in production", minProductionJWTSecretLen)
	}

	return nil
}
