// Package modelosenvase implements the modelos_envase CRUD usecases:
// orchestration between repository/modelos_envase and repository/auditoria.
// Every exported method returns either nil or a *domainerrors.DomainError.
package modelosenvase

import (
	"context"
	"encoding/json"
	"log/slog"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
)

const modelosEnvaseTable = "modelos_envase"

// Service groups every modelos_envase usecase behind one set of
// dependencies.
type Service struct {
	ModelosEnvase repo.Repository
	Auditoria     auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(modelosEnvaseRepo repo.Repository, auditoriaRepo auditoria.Repository) *Service {
	return &Service{ModelosEnvase: modelosEnvaseRepo, Auditoria: auditoriaRepo}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues.
type auditSnapshot struct {
	ID                 int64  `json:"id"`
	Tipo               string `json:"tipo"`
	TamanoOz           string `json:"tamano_oz"`
	EquivGramos        string `json:"equiv_gramos"`
	PrecioSolo         string `json:"precio_solo"`
	PrecioConFragancia string `json:"precio_con_fragancia"`
	PrecioRecarga      string `json:"precio_recarga"`
	Activo             bool   `json:"activo"`
}

func snapshot(m generated.ModelosEnvase) auditSnapshot {
	return auditSnapshot{
		ID:                 m.ID,
		Tipo:               m.Tipo,
		TamanoOz:           m.TamanoOz.String(),
		EquivGramos:        m.EquivGramos.String(),
		PrecioSolo:         m.PrecioSolo.String(),
		PrecioConFragancia: m.PrecioConFragancia.String(),
		PrecioRecarga:      m.PrecioRecarga.String(),
		Activo:             m.Activo,
	}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(modelosEnvaseTable),
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
	return domainerrors.NewNotFound("Modelo de envase no encontrado", "El modelo de envase solicitado no existe.")
}
