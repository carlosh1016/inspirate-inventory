package modelosenvase

import (
	"context"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// CreateInput is the request payload plus the requester's context.
type CreateInput struct {
	Tipo               string
	TamanoOz           string
	EquivGramos        string
	PrecioSolo         string
	PrecioConFragancia string
	PrecioRecarga      string
	RequesterID        int64
	IP                 string
	UserAgent          string
}

// Create registers a new modelo_envase, enforcing (tipo, tamano_oz)
// uniqueness.
func (s *Service) Create(ctx context.Context, in CreateInput) (generated.GetModeloEnvaseByIDRow, error) {
	exists, err := s.ModelosEnvase.ExistsTipoTamano(ctx, in.Tipo, in.TamanoOz, 0)
	if err != nil {
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}
	if exists {
		return generated.GetModeloEnvaseByIDRow{}, domainerrors.NewConflict(
			"Modelo en uso",
			"Ya existe un modelo de envase con ese tipo y tamaño.",
		)
	}

	m, err := s.ModelosEnvase.Insert(ctx, in.Tipo, in.TamanoOz, in.EquivGramos, in.PrecioSolo, in.PrecioConFragancia, in.PrecioRecarga)
	if err != nil {
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}

	result, err := s.ModelosEnvase.GetByID(ctx, m.ID)
	if err != nil {
		return generated.GetModeloEnvaseByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "modelo_envase_creado", in.IP, in.UserAgent, &m.ID, nil, snapshot(m))

	return result, nil
}
