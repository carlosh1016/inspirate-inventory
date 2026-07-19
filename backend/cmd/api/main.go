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
	fraganciashandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/fragancias"
	metodospagohandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/metodos_pago"
	modelosenvasehandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/modelos_envase"
	movimientoshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/movimientos"
	productoshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/productos"
	stockhandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/stock"
	usuarioshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/usuarios"
	variantesenvasehandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/variantes_envase"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/config"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/db"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/logger"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/mailer"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/ratelimit"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	metodospago "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
	modelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	movimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	passwordresets "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/password_resets"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	variantesenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	usecaseauth "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
	usecasefragancias "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
	usecasemetodospago "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/metodos_pago"
	usecasemodelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
	usecasemovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
	usecaseproductos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/productos"
	usecasestock "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/stock"
	usecaseusuarios "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
	usecasevariantesenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/variantes_envase"
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

	authHandler, usuariosHandler, fraganciasHandler, modelosEnvaseHandler, variantesEnvaseHandler, productosHandler, metodosPagoHandler, stockHandler, movimientosHandler := buildHandlers(cfg, pool)

	router := apphttp.NewRouter(cfg, log, pool, authHandler, usuariosHandler, fraganciasHandler, modelosEnvaseHandler, variantesEnvaseHandler, productosHandler, metodosPagoHandler, stockHandler, movimientosHandler)
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

func buildHandlers(cfg *config.Config, pool *pgxpool.Pool) (*authhandlers.Handler, *usuarioshandlers.Handler, *fraganciashandlers.Handler, *modelosenvasehandlers.Handler, *variantesenvasehandlers.Handler, *productoshandlers.Handler, *metodospagohandlers.Handler, *stockhandlers.Handler, *movimientoshandlers.Handler) {
	usuariosRepo := usuarios.NewPostgres(pool)
	refreshTokensRepo := refreshtokens.NewPostgres(pool)
	passwordResetsRepo := passwordresets.NewPostgres(pool)
	auditoriaRepo := auditoria.NewPostgres(pool)
	fraganciasRepo := fragancias.NewPostgres(pool)
	stockActualRepo := stockactual.NewPostgres(pool)
	modelosEnvaseRepo := modelosenvase.NewPostgres(pool)
	variantesEnvaseRepo := variantesenvase.NewPostgres(pool)
	productosRepo := productos.NewPostgres(pool)
	metodosPagoRepo := metodospago.NewPostgres(pool)
	movimientosRepo := movimientos.NewPostgres(pool)

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

	fraganciasService := usecasefragancias.NewService(pool, fraganciasRepo, stockActualRepo, auditoriaRepo)
	fraganciasHandler := fraganciashandlers.NewHandler(fraganciasService, jwtManager, v)

	modelosEnvaseService := usecasemodelosenvase.NewService(modelosEnvaseRepo, auditoriaRepo)
	modelosEnvaseHandler := modelosenvasehandlers.NewHandler(modelosEnvaseService, jwtManager, v)

	variantesEnvaseService := usecasevariantesenvase.NewService(pool, variantesEnvaseRepo, modelosEnvaseRepo, stockActualRepo, auditoriaRepo)
	variantesEnvaseHandler := variantesenvasehandlers.NewHandler(variantesEnvaseService, jwtManager, v)

	productosService := usecaseproductos.NewService(pool, productosRepo, stockActualRepo, auditoriaRepo)
	productosHandler := productoshandlers.NewHandler(productosService, jwtManager, v)

	metodosPagoService := usecasemetodospago.NewService(metodosPagoRepo, auditoriaRepo)
	metodosPagoHandler := metodospagohandlers.NewHandler(metodosPagoService, jwtManager, v)

	stockService := usecasestock.NewService(stockActualRepo)
	stockHandler := stockhandlers.NewHandler(stockService, jwtManager)

	movimientosService := usecasemovimientos.NewService(pool, movimientosRepo, stockActualRepo, fraganciasRepo, variantesEnvaseRepo, productosRepo, auditoriaRepo)
	movimientosHandler := movimientoshandlers.NewHandler(movimientosService, jwtManager, v)

	return authHandler, usuariosHandler, fraganciasHandler, modelosEnvaseHandler, variantesEnvaseHandler, productosHandler, metodosPagoHandler, stockHandler, movimientosHandler
}
