package fragancias

import (
	"context"
	"errors"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// UpdateInput is the request payload plus the requester's context. A nil
// field means "leave unchanged".
type UpdateInput struct {
	TargetID          int64
	NombreComercial   *string
	NombreAlternativo *string
	Genero            *string
	GramosMinimo      *string
	RequesterID       int64
	IP                string
	UserAgent         string
}

// Update applies a partial update to a fragancia, enforcing nombre_comercial
// uniqueness per sede when it changes.
func (s *Service) Update(ctx context.Context, in UpdateInput) (generated.GetFraganciaByIDRow, error) {
	before, err := s.Fragancias.GetByIDIncludingDeleted(ctx, in.TargetID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetFraganciaByIDRow{}, notFoundErr()
		}
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}
	if before.DeletedAt.Valid {
		return generated.GetFraganciaByIDRow{}, notFoundErr()
	}

	if in.NombreComercial != nil && *in.NombreComercial != before.NombreComercial {
		exists, err := s.Fragancias.ExistsNombreComercial(ctx, before.SedeID, *in.NombreComercial, in.TargetID)
		if err != nil {
			return generated.GetFraganciaByIDRow{}, internalErr(err)
		}
		if exists {
			return generated.GetFraganciaByIDRow{}, domainerrors.NewConflict(
				"Nombre en uso",
				"Ya existe una fragancia con ese nombre comercial en esta sede.",
			)
		}
	}

	updated, err := s.Fragancias.Update(ctx, in.TargetID, repo.UpdateFields{
		NombreComercial:   in.NombreComercial,
		NombreAlternativo: in.NombreAlternativo,
		Genero:            in.Genero,
		GramosMinimo:      in.GramosMinimo,
	})
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return generated.GetFraganciaByIDRow{}, notFoundErr()
		}
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}

	result, err := s.Fragancias.GetByID(ctx, updated.ID)
	if err != nil {
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "fragancia_editada", in.IP, in.UserAgent, &updated.ID, snapshot(before), snapshot(updated))

	return result, nil
}
