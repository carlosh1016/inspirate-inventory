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
	NumeroGenero      *int32
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

	if err := s.checkNumeroGeneroUnique(ctx, before, in); err != nil {
		return generated.GetFraganciaByIDRow{}, err
	}

	updated, err := s.Fragancias.Update(ctx, in.TargetID, repo.UpdateFields{
		NombreComercial:   in.NombreComercial,
		NombreAlternativo: in.NombreAlternativo,
		Genero:            in.Genero,
		GramosMinimo:      in.GramosMinimo,
		NumeroGenero:      in.NumeroGenero,
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

// checkNumeroGeneroUnique validates the (genero, numero_genero) pair when
// either is changing, comparing against the effective (post-update) values.
func (s *Service) checkNumeroGeneroUnique(ctx context.Context, before generated.Fragancia, in UpdateInput) error {
	if in.Genero == nil && in.NumeroGenero == nil {
		return nil
	}

	effectiveGenero := string(before.Genero)
	if in.Genero != nil {
		effectiveGenero = *in.Genero
	}
	effectiveNumero := before.NumeroGenero
	if in.NumeroGenero != nil {
		effectiveNumero = *in.NumeroGenero
	}
	if effectiveGenero == string(before.Genero) && effectiveNumero == before.NumeroGenero {
		return nil
	}

	exists, err := s.Fragancias.ExistsNumeroGenero(ctx, before.SedeID, effectiveGenero, effectiveNumero, in.TargetID)
	if err != nil {
		return internalErr(err)
	}
	if exists {
		return domainerrors.NewConflict("Número en uso", "Ya existe una fragancia con ese número en este género.")
	}
	return nil
}
