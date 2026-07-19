// Package fragancias implements the fragancias CRUD usecases: orchestration
// between domain/fragancias (pure rules), repository/fragancias,
// repository/stock_actual (to seed initial stock) and repository/auditoria.
// Every exported method returns either nil or a *domainerrors.DomainError.
package fragancias

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

const fraganciasTable = "fragancias"

// Service groups every fragancias usecase behind one set of dependencies.
// Pool is used only to open transactions for Create (insert + seed stock
// must be atomic); every other method goes through Fragancias/StockActual.
type Service struct {
	Pool        *pgxpool.Pool
	Fragancias  repo.Repository
	StockActual stockactual.Repository
	Auditoria   auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(pool *pgxpool.Pool, fraganciasRepo repo.Repository, stockActualRepo stockactual.Repository, auditoriaRepo auditoria.Repository) *Service {
	return &Service{Pool: pool, Fragancias: fraganciasRepo, StockActual: stockActualRepo, Auditoria: auditoriaRepo}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues.
type auditSnapshot struct {
	ID              int64  `json:"id"`
	NombreComercial string `json:"nombre_comercial"`
	Genero          string `json:"genero"`
	GramosMinimo    string `json:"gramos_minimo"`
	Activo          bool   `json:"activo"`
}

func snapshot(f generated.Fragancia) auditSnapshot {
	return auditSnapshot{
		ID:              f.ID,
		NombreComercial: f.NombreComercial,
		Genero:          string(f.Genero),
		GramosMinimo:    f.GramosMinimo.String(),
		Activo:          f.Activo,
	}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(fraganciasTable),
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
	return domainerrors.NewNotFound("Fragancia no encontrada", "La fragancia solicitada no existe.")
}
