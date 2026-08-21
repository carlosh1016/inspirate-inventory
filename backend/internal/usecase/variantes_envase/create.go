package variantesenvase

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	modelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

// CreateInput is the request payload plus the requester's context. SedeID is
// inherited from the requester's own claims, not client input.
type CreateInput struct {
	SedeID         int64
	ModeloEnvaseID int64
	Color          string
	StockMinimo    int32
	RequesterID    int64
	IP             string
	UserAgent      string
}

// Create registers a new variante_envase and seeds its stock (vitrina +
// bodega, both zero) atomically. The parent modelo_envase must exist and be
// active.
func (s *Service) Create(ctx context.Context, in CreateInput) (generated.GetVarianteEnvaseByIDRow, error) {
	modelo, err := s.ModelosEnvase.GetByID(ctx, in.ModeloEnvaseID)
	if err != nil {
		if errors.Is(err, modelosenvase.ErrNotFound) {
			return generated.GetVarianteEnvaseByIDRow{}, domainerrors.NewValidation(
				"Modelo de envase inválido",
				"El modelo de envase indicado no existe o fue eliminado.",
				nil,
			)
		}
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}
	if !modelo.Activo {
		return generated.GetVarianteEnvaseByIDRow{}, domainerrors.NewBusinessRule(
			"Modelo de envase inactivo",
			"No se pueden crear variantes para un modelo de envase inactivo.",
		)
	}
	if !modelo.TieneVariantes {
		return generated.GetVarianteEnvaseByIDRow{}, domainerrors.NewBusinessRule(
			"Modelo sin variantes",
			"Este modelo de envase es único y no admite variantes adicionales.",
		)
	}

	exists, err := s.VariantesEnvase.ExistsColor(ctx, in.ModeloEnvaseID, in.Color, 0)
	if err != nil {
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}
	if exists {
		return generated.GetVarianteEnvaseByIDRow{}, domainerrors.NewConflict(
			"Grosor en uso",
			"Ya existe una variante con ese grosor para este modelo de envase.",
		)
	}

	var variante generated.VariantesEnvase
	err = commonrepo.WithTx(ctx, s.Pool, func(tx pgx.Tx) error {
		txVariantes := repo.NewPostgres(tx)
		txStock := stockactual.NewPostgres(tx)

		v, err := txVariantes.Insert(ctx, in.ModeloEnvaseID, in.SedeID, in.Color, in.StockMinimo)
		if err != nil {
			return err
		}
		variante = v

		return txStock.InitializeStock(ctx, in.SedeID, stockactual.TipoItemVarianteEnvase, v.ID)
	})
	if err != nil {
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}

	result, err := s.VariantesEnvase.GetByID(ctx, variante.ID)
	if err != nil {
		return generated.GetVarianteEnvaseByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "variante_envase_creada", in.IP, in.UserAgent, &variante.ID, nil, snapshot(variante))

	return result, nil
}
