// Package movimientos implements the movimientos de inventario usecases:
// InventoryService (the transactional engine that applies stock changes
// safely under concurrency) plus the six operation-specific usecases built
// on top of it (entrada, traslado, ajuste, dañado, corrección, list).
// Every exported method returns either nil or a *domainerrors.DomainError.
package movimientos

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	fraganciasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	productosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactualrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	variantesenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

const movimientosTable = "movimientos_inventario"

// Service groups InventoryService and every movimientos usecase behind one
// set of dependencies. It implements InventoryService itself (see
// inventory_service.go).
type Service struct {
	Pool            *pgxpool.Pool
	Movimientos     movimientosrepo.Repository
	StockActual     stockactualrepo.Repository
	Fragancias      fraganciasrepo.Repository
	VariantesEnvase variantesenvaserepo.Repository
	Productos       productosrepo.Repository
	Auditoria       auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(
	pool *pgxpool.Pool,
	movimientosRepo movimientosrepo.Repository,
	stockActualRepo stockactualrepo.Repository,
	fraganciasRepo fraganciasrepo.Repository,
	variantesEnvaseRepo variantesenvaserepo.Repository,
	productosRepo productosrepo.Repository,
	auditoriaRepo auditoria.Repository,
) *Service {
	return &Service{
		Pool:            pool,
		Movimientos:     movimientosRepo,
		StockActual:     stockActualRepo,
		Fragancias:      fraganciasRepo,
		VariantesEnvase: variantesEnvaseRepo,
		Productos:       productosRepo,
		Auditoria:       auditoriaRepo,
	}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(movimientosTable),
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

// parsePositiveDecimal parses a decimal-string request field and rejects
// anything that isn't strictly greater than zero — used by every usecase
// that takes a plain "how much" quantity (entrada, traslado, dañado).
func parsePositiveDecimal(s string) (decimal.Decimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil || !d.IsPositive() {
		return decimal.Decimal{}, domainerrors.NewValidation(
			"Cantidad inválida",
			"La cantidad debe ser un número mayor que cero.",
			nil,
		)
	}
	return d, nil
}

// itemNoEncontradoErr reports that a movimiento references a catalog item
// that doesn't exist or is soft-deleted. Mapped to a validation error (422)
// rather than a bespoke 400 — this API has no 400 status, every input
// problem already goes through CodeValidation.
func itemNoEncontradoErr(tipoItem domainstock.TipoItem, itemID int64) error {
	return domainerrors.NewValidation(
		"Ítem inválido",
		fmt.Sprintf("El ítem %s con id %d no existe o fue eliminado.", tipoItem, itemID),
		nil,
	)
}

// toDomainMovimiento maps a generated row to the pure domain type. nombre
// is threaded through from the caller's own item-existence lookup — no
// extra query needed here.
func toDomainMovimiento(row generated.MovimientosInventario, nombre string) domainmovimientos.Movimiento {
	return domainmovimientos.Movimiento{
		ID:             row.ID,
		SedeID:         row.SedeID,
		UsuarioID:      row.UsuarioID,
		TipoItem:       domainstock.TipoItem(row.TipoItem),
		ItemID:         row.ItemID,
		ItemNombre:     nombre,
		Tipo:           domainmovimientos.TipoMovimiento(row.Tipo),
		Ubicacion:      domainstock.Ubicacion(row.Ubicacion),
		Cantidad:       row.Cantidad,
		StockAnterior:  row.StockAnterior,
		StockPosterior: row.StockPosterior,
		Motivo:         repo.StringPtr(row.Motivo),
		VentaID:        repo.Int8Ptr(row.VentaID),
		CreatedAt:      row.CreatedAt.Time,
	}
}
