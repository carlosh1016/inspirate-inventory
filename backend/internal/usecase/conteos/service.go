// Package conteos implements the monthly fragrance inventory count (conteo
// mensual de inventario) usecases: Crear, GetByID, List,
// ActualizarGramosFisico, Cerrar. Every exported method returns either nil
// or a *domainerrors.DomainError.
package conteos

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
)

const conteosTable = "conteos_inventario"

// Service groups every conteo mensual de inventario usecase behind one set
// of dependencies.
type Service struct {
	Pool      *pgxpool.Pool
	Conteos   conteosrepo.Repository
	Auditoria auditoria.Repository
	Location  *time.Location
}

// NewService builds a Service with all its dependencies. loc is the
// timezone used to resolve "el mes actual" when Periodo isn't given
// explicitly (America/Bogota, resolved by the caller — see cmd/api/main.go).
func NewService(pool *pgxpool.Pool, conteosRepo conteosrepo.Repository, auditoriaRepo auditoria.Repository, loc *time.Location) *Service {
	return &Service{
		Pool:      pool,
		Conteos:   conteosRepo,
		Auditoria: auditoriaRepo,
		Location:  loc,
	}
}

// primerDiaMesActual returns midnight-in-Colombia for the first day of the
// current month.
func (s *Service) primerDiaMesActual() time.Time {
	now := time.Now().In(s.Location)
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, s.Location)
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(conteosTable),
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
	return domainerrors.NewNotFound("Conteo no encontrado", "El conteo de inventario solicitado no existe.")
}

func itemNotFoundErr() error {
	return domainerrors.NewNotFound("Fragancia no encontrada en el conteo", "Esa fragancia no hace parte de este conteo.")
}
