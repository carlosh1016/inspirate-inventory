// Package cuadres implements the cuadre de caja (daily cash register
// close-out) usecases: Abrir, Cerrar, GetHoy, GetByID, List, AddPago,
// DeletePago, AddConsignacion, DeleteConsignacion — plus TotalesService
// (live ventas totals by payment method) and CajaStatusService (consumed by
// usecase/ventas to block sales once the day's cuadre is closed). Every
// exported method returns either nil or a *domainerrors.DomainError.
package cuadres

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	consignacionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/consignaciones"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	pagoscajarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/pagos_caja"
	usuariosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

const cuadresTable = "cuadres_caja"

// Service groups every cuadre de caja usecase behind one set of
// dependencies.
type Service struct {
	Pool           *pgxpool.Pool
	Cuadres        cuadresrepo.Repository
	Pagos          pagoscajarepo.Repository
	Consignaciones consignacionesrepo.Repository
	Usuarios       usuariosrepo.Repository
	Auditoria      auditoria.Repository
	Totales        TotalesService
	Location       *time.Location
}

// NewService builds a Service with all its dependencies. loc is the
// timezone used to compute "hoy" (America/Bogota, with a fixed-offset
// fallback resolved by the caller — see cmd/api/main.go).
func NewService(
	pool *pgxpool.Pool,
	cuadresRepo cuadresrepo.Repository,
	pagosRepo pagoscajarepo.Repository,
	consignacionesRepo consignacionesrepo.Repository,
	usuariosRepo usuariosrepo.Repository,
	auditoriaRepo auditoria.Repository,
	loc *time.Location,
) *Service {
	return &Service{
		Pool:           pool,
		Cuadres:        cuadresRepo,
		Pagos:          pagosRepo,
		Consignaciones: consignacionesRepo,
		Usuarios:       usuariosRepo,
		Auditoria:      auditoriaRepo,
		Totales:        NewTotalesService(cuadresRepo, loc),
		Location:       loc,
	}
}

// hoy returns midnight-in-Colombia for "today", the same computation used
// consistently across abrir/get-hoy/caja-status.
func (s *Service) hoy() time.Time {
	now := time.Now().In(s.Location)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.Location)
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(cuadresTable),
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
	return domainerrors.NewNotFound("Cuadre no encontrado", "El cuadre de caja solicitado no existe.")
}

func pagoNotFoundErr() error {
	return domainerrors.NewNotFound("Pago no encontrado", "El pago solicitado no existe.")
}

func consignacionNotFoundErr() error {
	return domainerrors.NewNotFound("Consignación no encontrada", "La consignación solicitada no existe.")
}
