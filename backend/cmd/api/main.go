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

	_ "time/tzdata" // ensures America/Bogota resolves without relying on OS tzdata

	"github.com/jackc/pgx/v5/pgxpool"

	apphttp "github.com/carlosh1016/inspirate-inventory/backend/internal/http"
	auditoriahandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/auditoria"
	authhandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/auth"
	cuadreshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/cuadres"
	fraganciashandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/fragancias"
	metodospagohandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/metodos_pago"
	modelosenvasehandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/modelos_envase"
	movimientoshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/movimientos"
	productoshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/productos"
	reporteshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/reportes"
	sesioneshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/sesiones"
	stockhandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/stock"
	usuarioshandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/usuarios"
	variantesenvasehandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/variantes_envase"
	ventashandlers "github.com/carlosh1016/inspirate-inventory/backend/internal/http/handlers/ventas"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/config"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/db"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/logger"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/mailer"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/ratelimit"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/consignaciones"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	idempotencykeys "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/idempotency_keys"
	metodospago "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
	modelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	movimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	pagoscaja "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/pagos_caja"
	passwordresets "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/password_resets"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	reportesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/reportes"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	variantesenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	ventaitems "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/venta_items"
	ventasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/ventas"
	usecaseauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auditoria"
	usecaseauth "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
	usecasecuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
	usecasefragancias "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
	usecasemetodospago "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/metodos_pago"
	usecasemodelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
	usecasemovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
	usecaseproductos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/productos"
	usecasereportes "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/reportes"
	usecasesesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
	usecasestock "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/stock"
	usecaseusuarios "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
	usecasevariantesenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/variantes_envase"
	usecaseventas "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
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

	authHandler, usuariosHandler, fraganciasHandler, modelosEnvaseHandler, variantesEnvaseHandler, productosHandler, metodosPagoHandler, stockHandler, movimientosHandler, ventasHandler, cuadresHandler, sesionesHandler, reportesHandler, auditoriaHandler, idempotencyKeysRepo := buildHandlers(cfg, pool, log)

	router := apphttp.NewRouter(cfg, log, pool, authHandler, usuariosHandler, fraganciasHandler, modelosEnvaseHandler, variantesEnvaseHandler, productosHandler, metodosPagoHandler, stockHandler, movimientosHandler, ventasHandler, cuadresHandler, sesionesHandler, reportesHandler, auditoriaHandler)
	server := apphttp.NewServer(cfg.Port, router)

	go runIdempotencyKeyCleanup(ctx, idempotencyKeysRepo, log)

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

// idempotencyCleanupInterval is how often expired idempotency_keys rows are
// purged. An hour is frequent enough given the 24h TTL those rows carry,
// and infrequent enough not to matter for load.
const idempotencyCleanupInterval = 1 * time.Hour

// runIdempotencyKeyCleanup periodically deletes expired idempotency_keys
// rows until ctx is done (the same context signal.NotifyContext ties to
// shutdown) — fire-and-forget from run()'s perspective, no separate cancel
// plumbing needed.
func runIdempotencyKeyCleanup(ctx context.Context, repo idempotencykeys.Repository, log *slog.Logger) {
	ticker := time.NewTicker(idempotencyCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := repo.DeleteExpired(ctx); err != nil {
				log.ErrorContext(ctx, "failed to delete expired idempotency keys", "error", err)
			}
		}
	}
}

// resolveBogotaLocation loads America/Bogota, falling back to a fixed
// UTC-5 offset if the timezone database is unavailable — Colombia has no
// daylight saving, so UTC-5 is exact year-round. The blank import of
// time/tzdata above means this fallback essentially never triggers in
// practice; it exists as defense in depth, not as the expected path.
func resolveBogotaLocation(log *slog.Logger) *time.Location {
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		log.Warn("could not load America/Bogota timezone, falling back to fixed UTC-5", "error", err)
		return time.FixedZone("America/Bogota", -5*60*60)
	}
	return loc
}

func buildHandlers(cfg *config.Config, pool *pgxpool.Pool, log *slog.Logger) (*authhandlers.Handler, *usuarioshandlers.Handler, *fraganciashandlers.Handler, *modelosenvasehandlers.Handler, *variantesenvasehandlers.Handler, *productoshandlers.Handler, *metodospagohandlers.Handler, *stockhandlers.Handler, *movimientoshandlers.Handler, *ventashandlers.Handler, *cuadreshandlers.Handler, *sesioneshandlers.Handler, *reporteshandlers.Handler, *auditoriahandlers.Handler, idempotencykeys.Repository) {
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
	ventasRepo := ventasrepo.NewPostgres(pool)
	ventaItemsRepo := ventaitems.NewPostgres(pool)
	idempotencyKeysRepo := idempotencykeys.NewPostgres(pool)
	cuadresRepo := cuadresrepo.NewPostgres(pool)
	pagosCajaRepo := pagoscaja.NewPostgres(pool)
	consignacionesRepo := consignaciones.NewPostgres(pool)
	sesionesRepo := sesionesrepo.NewPostgres(pool)
	reportesRepo := reportesrepo.NewPostgres(pool)

	jwtManager := jwt.New(cfg.JWTSecret)
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
		cfg.JWTAccessTTLAdmin,
		cfg.JWTAccessTTLVendedora,
		cfg.JWTRefreshTTLAdmin,
		cfg.JWTRefreshTTLVendedora,
		cfg.FrontendURL,
	)
	authHandler := authhandlers.NewHandler(authService, jwtManager, v, cfg.Environment == "production")

	usuariosService := usecaseusuarios.NewService(usuariosRepo, refreshTokensRepo, auditoriaRepo)
	usuariosHandler := usuarioshandlers.NewHandler(usuariosService, jwtManager, v)

	fraganciasService := usecasefragancias.NewService(pool, fraganciasRepo, stockActualRepo, auditoriaRepo)
	fraganciasHandler := fraganciashandlers.NewHandler(fraganciasService, jwtManager, v)

	modelosEnvaseService := usecasemodelosenvase.NewService(pool, modelosEnvaseRepo, variantesEnvaseRepo, stockActualRepo, auditoriaRepo)
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

	ventasService := usecaseventas.NewService(
		pool,
		ventasRepo,
		ventaItemsRepo,
		idempotencyKeysRepo,
		movimientosService,
		movimientosRepo,
		fraganciasRepo,
		variantesEnvaseRepo,
		modelosEnvaseRepo,
		productosRepo,
		metodosPagoRepo,
		auditoriaRepo,
		usecaseventas.NewPricingService(),
		usecaseventas.NewDiscountService(),
		usecasecuadres.NewCajaStatusService(cuadresRepo),
		resolveBogotaLocation(log),
	)
	ventasHandler := ventashandlers.NewHandler(ventasService, jwtManager, v)

	cuadresService := usecasecuadres.NewService(pool, cuadresRepo, pagosCajaRepo, consignacionesRepo, usuariosRepo, auditoriaRepo, resolveBogotaLocation(log))
	cuadresHandler := cuadreshandlers.NewHandler(cuadresService, jwtManager, v)

	sesionesService := usecasesesiones.NewService(sesionesRepo, usuariosRepo, auditoriaRepo, resolveBogotaLocation(log))
	sesionesHandler := sesioneshandlers.NewHandler(sesionesService, jwtManager, v)

	reportesService := usecasereportes.NewService(reportesRepo, resolveBogotaLocation(log))
	reportesHandler := reporteshandlers.NewHandler(reportesService, jwtManager, resolveBogotaLocation(log))

	auditoriaService := usecaseauditoria.NewService(auditoriaRepo)
	auditoriaHandler := auditoriahandlers.NewHandler(auditoriaService, jwtManager)

	return authHandler, usuariosHandler, fraganciasHandler, modelosEnvaseHandler, variantesEnvaseHandler, productosHandler, metodosPagoHandler, stockHandler, movimientosHandler, ventasHandler, cuadresHandler, sesionesHandler, reportesHandler, auditoriaHandler, idempotencyKeysRepo
}
