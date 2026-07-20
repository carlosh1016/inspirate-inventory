// Package ventas implements the ventas usecases: CreateVenta (the complex
// one — pricing, discount, stock consolidation, idempotency, all inside one
// transaction that also calls into usecase/movimientos), plus List, Get,
// Update (observaciones only) and ResumenHoy. Every exported method returns
// either nil or a *domainerrors.DomainError.
package ventas

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	fraganciasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	idempotencykeysrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/idempotency_keys"
	metodospagorepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
	modelosenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	productosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	variantesenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	ventaitemsrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/venta_items"
	ventasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/ventas"
	usecasemovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

const ventasTable = "ventas"

// Service groups every ventas usecase behind one set of dependencies. Pool
// is used to open the RepeatableRead transaction CreateVenta needs to load
// entities, insert the venta+items, and call Movimientos.RegisterBatchTx
// all atomically.
type Service struct {
	Pool            *pgxpool.Pool
	Ventas          ventasrepo.Repository
	VentaItems      ventaitemsrepo.Repository
	IdempotencyKeys idempotencykeysrepo.Repository
	Movimientos     usecasemovimientos.InventoryService
	MovimientosRepo movimientosrepo.Repository
	Fragancias      fraganciasrepo.Repository
	VariantesEnvase variantesenvaserepo.Repository
	ModelosEnvase   modelosenvaserepo.Repository
	Productos       productosrepo.Repository
	MetodosPago     metodospagorepo.Repository
	Auditoria       auditoria.Repository
	Pricing         domainventas.PricingService
	Discount        domainventas.DiscountService
	Location        *time.Location
}

// NewService builds a Service with all its dependencies. loc is the
// timezone used by ResumenHoy to compute "today" (America/Bogota, with a
// fixed-offset fallback resolved by the caller — see cmd/api/main.go).
func NewService(
	pool *pgxpool.Pool,
	ventasRepo ventasrepo.Repository,
	ventaItemsRepo ventaitemsrepo.Repository,
	idempotencyKeysRepo idempotencykeysrepo.Repository,
	movimientosService usecasemovimientos.InventoryService,
	movimientosRepo movimientosrepo.Repository,
	fraganciasRepo fraganciasrepo.Repository,
	variantesEnvaseRepo variantesenvaserepo.Repository,
	modelosEnvaseRepo modelosenvaserepo.Repository,
	productosRepo productosrepo.Repository,
	metodosPagoRepo metodospagorepo.Repository,
	auditoriaRepo auditoria.Repository,
	pricing domainventas.PricingService,
	discount domainventas.DiscountService,
	loc *time.Location,
) *Service {
	return &Service{
		Pool:            pool,
		Ventas:          ventasRepo,
		VentaItems:      ventaItemsRepo,
		IdempotencyKeys: idempotencyKeysRepo,
		Movimientos:     movimientosService,
		MovimientosRepo: movimientosRepo,
		Fragancias:      fraganciasRepo,
		VariantesEnvase: variantesEnvaseRepo,
		ModelosEnvase:   modelosEnvaseRepo,
		Productos:       productosRepo,
		MetodosPago:     metodosPagoRepo,
		Auditoria:       auditoriaRepo,
		Pricing:         pricing,
		Discount:        discount,
		Location:        loc,
	}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues
// for venta_observaciones_editadas.
type auditSnapshot struct {
	Observaciones *string `json:"observaciones"`
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(ventasTable),
		RegistroID:    registroID,
		IP:            ip,
		UserAgent:     userAgent,
	}
	if antes != nil {
		if b, err := json.Marshal(antes); err == nil {
			entry.DatosAntes = b
		}
	}
	if despues != nil {
		if b, err := json.Marshal(despues); err == nil {
			entry.DatosDespues = b
		}
	}

	if err := s.Auditoria.Insert(ctx, entry); err != nil {
		slog.ErrorContext(ctx, "failed to write auditoria entry", "accion", accion, "error", err)
	}
}

func strPtr(s string) *string { return &s }

func internalErr(err error) error {
	return domainerrors.NewInternal("Error interno", "Ocurrió un error inesperado. Intenta de nuevo más tarde.", err)
}

func notFoundErr() error {
	return domainerrors.NewNotFound("Venta no encontrada", "La venta solicitada no existe.")
}
