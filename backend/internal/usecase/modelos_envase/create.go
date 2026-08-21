package modelosenvase

import (
	"context"

	"github.com/jackc/pgx/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	variantesenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
)

// CreateInput is the request payload plus the requester's context. SedeID is
// inherited from the requester's own claims — it's only used when
// SinVariantes is true, to seed the hidden variante's stock for the
// requester's sede. SinVariantes defaults to false (the common case: a
// modelo that varies by grosor) so callers that don't set it get the normal
// behavior, not the "envase de lujo" one.
type CreateInput struct {
	Tipo               string
	TamanoOz           string
	EquivGramos        string
	PrecioSolo         string
	PrecioConFragancia string
	PrecioRecarga      string
	SinVariantes       bool
	SedeID             int64
	RequesterID        int64
	IP                 string
	UserAgent          string
}

// Create registers a new modelo_envase, enforcing (tipo, tamano_oz)
// uniqueness. When SinVariantes is true (e.g. an "envase de lujo" that never
// varies by grosor/color), it also creates a single hidden variante_envase
// for it, with its stock seeded, atomically — the UI never exposes this
// variante for manual management.
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

	var m generated.ModelosEnvase
	err = commonrepo.WithTx(ctx, s.Pool, func(tx pgx.Tx) error {
		txModelos := repo.NewPostgres(tx)

		inserted, err := txModelos.Insert(ctx, in.Tipo, in.TamanoOz, in.EquivGramos, in.PrecioSolo, in.PrecioConFragancia, in.PrecioRecarga, !in.SinVariantes)
		if err != nil {
			return err
		}
		m = inserted

		if !in.SinVariantes {
			return nil
		}

		txVariantes := variantesenvase.NewPostgres(tx)
		txStock := stockactual.NewPostgres(tx)

		v, err := txVariantes.Insert(ctx, m.ID, in.SedeID, varianteUnicaGrosor, 0)
		if err != nil {
			return err
		}

		return txStock.InitializeStock(ctx, in.SedeID, stockactual.TipoItemVarianteEnvase, v.ID)
	})
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
