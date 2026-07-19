// Package http contains the central route registration and router
// configuration.
package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers"
	authhandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/auth"
	usuarioshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/usuarios"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/config"
)

// NewRouter builds the chi router with base middleware and every registered
// route group.
func NewRouter(
	cfg *config.Config,
	logger *slog.Logger,
	pool *pgxpool.Pool,
	authHandler *authhandlers.Handler,
	usuariosHandler *usuarioshandlers.Handler,
) http.Handler {
	logger.Debug("building router", "cors_origins", cfg.CORSAllowedOrigins)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.AuditContext)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handlers.Health(pool))
		authHandler.Router(r)
		usuariosHandler.Router(r)
	})

	return r
}
