// Command api starts the Inspírate Inventory HTTP server.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	apphttp "github.com/carlosh1016/inspirate-inventory/backend/internal/http"
	authhandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/auth"
	usuarioshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/usuarios"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/config"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/db"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/logger"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/mailer"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/ratelimit"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	passwordresets "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/password_resets"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	usecaseauth "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
	usecaseusuarios "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg.Environment, cfg.LogLevel)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL, log)
	if err != nil {
		return err
	}
	defer pool.Close()

	authHandler, usuariosHandler := buildHandlers(cfg, pool)

	router := apphttp.NewRouter(cfg, log, pool, authHandler, usuariosHandler)
	server := apphttp.NewServer(cfg.Port, router)

	serverErr := make(chan error, 1)
	go func() {
		log.Info("starting server", "port", cfg.Port, "environment", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Info("server stopped gracefully")
	return nil
}

func buildHandlers(cfg *config.Config, pool *pgxpool.Pool) (*authhandlers.Handler, *usuarioshandlers.Handler) {
	usuariosRepo := usuarios.NewPostgres(pool)
	refreshTokensRepo := refreshtokens.NewPostgres(pool)
	passwordResetsRepo := passwordresets.NewPostgres(pool)
	auditoriaRepo := auditoria.NewPostgres(pool)

	jwtManager := jwt.New(cfg.JWTSecret, cfg.JWTAccessTTL)
	v := validator.New()

	var mailerSvc mailer.Mailer
	if cfg.Environment == "production" && cfg.ResendAPIKey != "" {
		mailerSvc = mailer.NewResendMailer(cfg.ResendAPIKey, cfg.MailFrom)
	} else {
		mailerSvc = mailer.NewMock()
	}

	authService := usecaseauth.NewService(
		usuariosRepo,
		refreshTokensRepo,
		passwordResetsRepo,
		auditoriaRepo,
		jwtManager,
		mailerSvc,
		ratelimit.NewLoginLimiter(),
		ratelimit.NewPasswordResetLimiter(),
		cfg.JWTRefreshTTLAdmin,
		cfg.JWTRefreshTTLVendedora,
		cfg.FrontendURL,
	)
	authHandler := authhandlers.NewHandler(authService, jwtManager, v, cfg.Environment == "production")

	usuariosService := usecaseusuarios.NewService(usuariosRepo, refreshTokensRepo, auditoriaRepo)
	usuariosHandler := usuarioshandlers.NewHandler(usuariosService, jwtManager, v)

	return authHandler, usuariosHandler
}
