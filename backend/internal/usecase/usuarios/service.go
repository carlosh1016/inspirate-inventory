// Package usuarios implements the usuarios CRUD usecases: orchestration
// between domain/usuarios (pure rules), domain/auth (password hashing,
// reused rather than duplicated), and the repository layer. Every exported
// method returns either nil or a *domainerrors.DomainError.
package usuarios

import (
	"context"
	"encoding/json"
	"log/slog"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
)

const usuariosTable = "usuarios"

// Service groups every usuarios usecase behind one set of dependencies.
type Service struct {
	Usuarios      repo.Repository
	RefreshTokens refreshtokens.Repository
	Auditoria     auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(usuariosRepo repo.Repository, refreshTokensRepo refreshtokens.Repository, auditoriaRepo auditoria.Repository) *Service {
	return &Service{Usuarios: usuariosRepo, RefreshTokens: refreshTokensRepo, Auditoria: auditoriaRepo}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues —
// deliberately excludes password_hash.
type auditSnapshot struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
	Correo         string `json:"correo"`
	Rol            string `json:"rol"`
	IsActive       bool   `json:"is_active"`
}

func snapshot(u generated.Usuario) auditSnapshot {
	return auditSnapshot{
		ID:             u.ID,
		NombreCompleto: u.NombreCompleto,
		Correo:         u.Correo,
		Rol:            string(u.Rol),
		IsActive:       u.IsActive,
	}
}

// audit records an auditoria entry. Failures are logged, never returned —
// auditing must never block the primary action it's recording.
func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(usuariosTable),
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
	return domainerrors.NewNotFound("Usuario no encontrado", "El usuario solicitado no existe.")
}
