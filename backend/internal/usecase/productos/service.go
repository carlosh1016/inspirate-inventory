// Package productos implements the productos CRUD usecases: orchestration
// between domain/productos (pure rules), repository/productos,
// repository/stock_actual (to seed initial stock) and repository/auditoria.
// Every exported method returns either nil or a *domainerrors.DomainError.
package productos

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

const productosTable = "productos"

// Service groups every productos usecase behind one set of dependencies.
// Pool is used only to open transactions for Create (insert + seed stock
// must be atomic); every other method goes through Productos/StockActual.
type Service struct {
	Pool        *pgxpool.Pool
	Productos   repo.Repository
	StockActual stockactual.Repository
	Auditoria   auditoria.Repository
}

// NewService builds a Service with all its dependencies.
func NewService(pool *pgxpool.Pool, productosRepo repo.Repository, stockActualRepo stockactual.Repository, auditoriaRepo auditoria.Repository) *Service {
	return &Service{Pool: pool, Productos: productosRepo, StockActual: stockActualRepo, Auditoria: auditoriaRepo}
}

// auditSnapshot is what gets JSON-encoded into datos_antes/datos_despues.
type auditSnapshot struct {
	ID          int64  `json:"id"`
	Nombre      string `json:"nombre"`
	Categoria   string `json:"categoria"`
	Precio      string `json:"precio"`
	StockMinimo int32  `json:"stock_minimo"`
	Activo      bool   `json:"activo"`
}

func snapshot(p generated.Producto) auditSnapshot {
	return auditSnapshot{
		ID:          p.ID,
		Nombre:      p.Nombre,
		Categoria:   string(p.Categoria),
		Precio:      p.Precio.String(),
		StockMinimo: p.StockMinimo,
		Activo:      p.Activo,
	}
}

func (s *Service) audit(ctx context.Context, requesterID *int64, accion, ip, userAgent string, registroID *int64, antes, despues any) {
	entry := auditoria.Entry{
		UsuarioID:     requesterID,
		Accion:        accion,
		TablaAfectada: strPtr(productosTable),
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
	return domainerrors.NewNotFound("Producto no encontrado", "El producto solicitado no existe.")
}
