package sesiones

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainsesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/sesiones"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
)

// UpdateInput is the request payload plus the requester's context. nil
// fields keep the existing value (mirrors sqlc.narg's COALESCE at the
// repository layer).
type UpdateInput struct {
	TargetID    int64
	EntradaAt   *time.Time
	SalidaAt    *time.Time
	RequesterID int64
	IP          string
	UserAgent   string
}

// auditSesionSnapshot is what gets JSON-encoded into datos_antes/
// datos_despues for sesion_editada.
type auditSesionSnapshot struct {
	EntradaAt time.Time  `json:"entrada_at"`
	SalidaAt  *time.Time `json:"salida_at"`
}

// Update lets an admin manually correct entrada_at/salida_at.
// horas_trabajadas is recalculated automatically by the repository layer.
func (s *Service) Update(ctx context.Context, in UpdateInput) (*domainsesiones.Sesion, error) {
	before, err := s.Sesiones.GetByID(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, sesionesrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}

	entradaEfectiva := before.EntradaAt.Time
	if in.EntradaAt != nil {
		entradaEfectiva = *in.EntradaAt
	}
	salidaEfectiva := repo.TimePtr(before.SalidaAt)
	if in.SalidaAt != nil {
		salidaEfectiva = in.SalidaAt
	}
	if salidaEfectiva != nil && entradaEfectiva.After(*salidaEfectiva) {
		return nil, domainerrors.NewValidation("Fechas inválidas", "entrada_at no puede ser posterior a salida_at.", nil)
	}

	row, err := s.Sesiones.UpdateManual(ctx, sesionesrepo.UpdateManualParams{
		ID:        in.TargetID,
		EntradaAt: in.EntradaAt,
		SalidaAt:  in.SalidaAt,
	})
	if err != nil {
		return nil, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "sesion_editada", in.IP, in.UserAgent, &row.ID,
		auditSesionSnapshot{EntradaAt: before.EntradaAt.Time, SalidaAt: repo.TimePtr(before.SalidaAt)},
		auditSesionSnapshot{EntradaAt: row.EntradaAt.Time, SalidaAt: repo.TimePtr(row.SalidaAt)},
	)

	sesion := toDomainSesion(row)
	return &sesion, nil
}
