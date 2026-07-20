// Package sesiones implements the sesiones laborales (vendedora clock
// in/out) usecases: Entrada, Salida, List, Update (admin manual
// correction), Resumen. Every exported method returns either nil or a
// *domainerrors.DomainError.
package sesiones

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
	usuariosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

const sesionesTable = "sesiones_laborales"

// Service groups every sesiones laborales usecase behind one set of
// dependencies.
type Service struct {
	Sesiones  sesionesrepo.Repository
	Usuarios  usuariosrepo.Repository
	Auditoria auditoria.Repository
	Location  *time.Location
}

// NewService builds a Service with all its dependencies. loc is the
// timezone used by Resumen's dias_trabajados grouping (America/Bogota,
// with a fixed-offset fallback resolved by the caller — see cmd/api/main.go).
func NewService(sesionesRepo sesionesrepo.Repository, usuariosRepo usuariosrepo.Repository, auditoriaRepo auditoria.Repository, loc *time.Location) *Service {
	return &Service{
		Sesiones:  sesionesRepo,
		Usuarios:  usuariosRepo,
		Auditoria: auditoriaRepo,
		Location:  loc,
	}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(sesionesTable),
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
	return domainerrors.NewNotFound("Sesión no encontrada", "La sesión laboral solicitada no existe.")
}
