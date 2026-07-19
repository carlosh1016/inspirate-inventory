// Package metodospago implements the metodos_pago CRUD usecases:
// orchestration between repository/metodos_pago and repository/auditoria.
// Every exported method returns either nil or a *domainerrors.DomainError.
package metodospago

import (
	"context"
	"encoding/json"
	"log/slog"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
)

const metodosPagoTable = "metodos_pago"

// Service groups every metodos_pago usecase behind one set of
// dependencies.
type Service struct {
	MetodosPago repo.Repository
	Auditoria   auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(metodosPagoRepo repo.Repository, auditoriaRepo auditoria.Repository) *Service {
	return &Service{MetodosPago: metodosPagoRepo, Auditoria: auditoriaRepo}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues.
type auditSnapshot struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
	Codigo string `json:"codigo"`
	Activo bool   `json:"activo"`
}

func snapshot(m generated.MetodosPago) auditSnapshot {
	return auditSnapshot{ID: m.ID, Nombre: m.Nombre, Codigo: m.Codigo, Activo: m.Activo}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(metodosPagoTable),
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
	return domainerrors.NewNotFound("Método de pago no encontrado", "El método de pago solicitado no existe.")
}
