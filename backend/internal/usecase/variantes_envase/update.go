package variantesenvase

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

// UpdateInput is the request payload plus the requester's context. A nil
// field means "leave unchanged". Both admin and vendedora share this same
// set of editable fields — there's no field variantes_envase exposes only
// to admin on update.
type UpdateInput struct {
	TargetID    int64
	Color       *string
	StockMinimo *int32
	RequesterID int64
	IP          string
	UserAgent   string
}

// Update applies a partial update to a variante_envase, enforcing color
// uniqueness per modelo_envase when it changes.
func (s *Service) Update(ctx context.Context, in UpdateInput) (generated.GetVarianteEnvaseByIDRow, error) {
	before, err := s.VariantesEnvase.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetVarianteEnvaseByIDRow{}, notFoundErr()
		}
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}
	if before.DeletedAt.Valid {
		return generated.GetVarianteEnvaseByIDRow{}, notFoundErr()
	}

	if in.Color != nil && *in.Color != before.Color {
		exists, err := s.VariantesEnvase.ExistsColor(ctx, before.ModeloEnvaseID, *in.Color, in.TargetID)
		if err != nil {
			return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
		}
		if exists {
			return generated.GetVarianteEnvaseByIDRow{}, domainerrors.NewConflict(
				"Color en uso",
				"Ya existe una variante con ese color para este modelo de envase.",
			)
		}
	}

	updated, err := s.VariantesEnvase.Update(ctx, in.TargetID, repo.UpdateFields{
		Color:       in.Color,
		StockMinimo: in.StockMinimo,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetVarianteEnvaseByIDRow{}, notFoundErr()
		}
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}

	result, err := s.VariantesEnvase.GetByID(ctx, updated.ID)
	if err != nil {
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "variante_envase_editada", in.IP, in.UserAgent, &updated.ID, snapshot(before), snapshot(updated))

	return result, nil
}
