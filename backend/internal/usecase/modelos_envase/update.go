package modelosenvase

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
)

// UpdateInput is the request payload plus the requester's context. A nil
// field means "leave unchanged".
type UpdateInput struct {
	TargetID           int64
	Tipo               *string
	TamanoOz           *string
	EquivGramos        *string
	PrecioSolo         *string
	PrecioConFragancia *string
	PrecioRecarga      *string
	RequesterID        int64
	IP                 string
	UserAgent          string
}

// Update applies a partial update to a modelo_envase, enforcing (tipo,
// tamano_oz) uniqueness when either changes.
func (s *Service) Update(ctx context.Context, in UpdateInput) (generated.GetModeloEnvaseByIDRow, error) {
	before, err := s.ModelosEnvase.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetModeloEnvaseByIDRow{}, notFoundErr()
		}
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}
	if before.DeletedAt.Valid {
		return generated.GetModeloEnvaseByIDRow{}, notFoundErr()
	}

	tipo := before.Tipo
	if in.Tipo != nil {
		tipo = *in.Tipo
	}
	tamanoOz := before.TamanoOz.String()
	if in.TamanoOz != nil {
		tamanoOz = *in.TamanoOz
	}
	if in.Tipo != nil || in.TamanoOz != nil {
		exists, err := s.ModelosEnvase.ExistsTipoTamano(ctx, tipo, tamanoOz, in.TargetID)
		if err != nil {
			return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
		}
		if exists {
			return generated.GetModeloEnvaseByIDRow{}, domainerrors.NewConflict(
				"Modelo en uso",
				"Ya existe un modelo de envase con ese tipo y tamaño.",
			)
		}
	}

	updated, err := s.ModelosEnvase.Update(ctx, in.TargetID, repo.UpdateFields{
		Tipo:               in.Tipo,
		TamanoOz:           in.TamanoOz,
		EquivGramos:        in.EquivGramos,
		PrecioSolo:         in.PrecioSolo,
		PrecioConFragancia: in.PrecioConFragancia,
		PrecioRecarga:      in.PrecioRecarga,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetModeloEnvaseByIDRow{}, notFoundErr()
		}
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}

	result, err := s.ModelosEnvase.GetByID(ctx, updated.ID)
	if err != nil {
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "modelo_envase_editado", in.IP, in.UserAgent, &updated.ID, snapshot(before), snapshot(updated))

	return result, nil
}
