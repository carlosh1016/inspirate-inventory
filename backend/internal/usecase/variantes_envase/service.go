// Package variantesenvase implements the variantes_envase CRUD usecases:
// orchestration between repository/variantes_envase, repository/modelos_envase
// (to validate the parent modelo), repository/stock_actual (to seed initial
// stock) and repository/auditoria. Every exported method returns either nil
// or a *domainerrors.DomainError.
package variantesenvase

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	modelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

const variantesEnvaseTable = "variantes_envase"

// Service groups every variantes_envase usecase behind one set of
// dependencies. Pool is used only to open transactions for Create (insert +
// seed stock must be atomic); every other method goes through the repos
// directly.
type Service struct {
	Pool            *pgxpool.Pool
	VariantesEnvase repo.Repository
	ModelosEnvase   modelosenvase.Repository
	StockActual     stockactual.Repository
	Auditoria       auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(pool *pgxpool.Pool, variantesEnvaseRepo repo.Repository, modelosEnvaseRepo modelosenvase.Repository, stockActualRepo stockactual.Repository, auditoriaRepo auditoria.Repository) *Service {
	return &Service{
		Pool:            pool,
		VariantesEnvase: variantesEnvaseRepo,
		ModelosEnvase:   modelosEnvaseRepo,
		StockActual:     stockActualRepo,
		Auditoria:       auditoriaRepo,
	}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues.
type auditSnapshot struct {
	ID             int64  `json:"id"`
	ModeloEnvaseID int64  `json:"modelo_envase_id"`
	Color          string `json:"color"`
	StockMinimo    int32  `json:"stock_minimo"`
	Activo         bool   `json:"activo"`
}

func snapshot(v generated.VariantesEnvase) auditSnapshot {
	return auditSnapshot{
		ID:             v.ID,
		ModeloEnvaseID: v.ModeloEnvaseID,
		Color:          v.Color,
		StockMinimo:    v.StockMinimo,
		Activo:         v.Activo,
	}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(variantesEnvaseTable),
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
	return domainerrors.NewNotFound("Variante de envase no encontrada", "La variante de envase solicitada no existe.")
}
