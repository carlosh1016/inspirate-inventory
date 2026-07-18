// Package config carga la configuración de la aplicación desde variables de entorno.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación, cargada desde el entorno.
type Config struct {
	Port                   string        `env:"PORT" envDefault:"8080"`
	Environment            string        `env:"ENVIRONMENT" envDefault:"development"`
	LogLevel               string        `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseURL            string        `env:"DATABASE_URL"`
	FrontendURL            string        `env:"FRONTEND_URL" envDefault:"http://localhost:3000"`
	CORSOrigins            []string      `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:3000"`
	JWTSecret              string        `env:"JWT_SECRET"`
	JWTAccessTTL           time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	JWTRefreshTTLAdmin     time.Duration `env:"JWT_REFRESH_TTL_ADMIN" envDefault:"720h"`
	JWTRefreshTTLVendedora time.Duration `env:"JWT_REFRESH_TTL_VENDEDORA" envDefault:"8h"`
	ResendAPIKey           string        `env:"RESEND_API_KEY"`
	MailFrom               string        `env:"MAIL_FROM"`
}

// Load lee un archivo .env si existe (ignorando el error si no está presente,
// típico en producción donde las variables vienen del entorno del sistema) y
// luego parsea la configuración desde las variables de entorno.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}
