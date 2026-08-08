package fragancias

import (
	"context"

	"github.com/jackc/pgx/v5"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
)

// CreateInput is the request payload plus the requester's context. SedeID is
// inherited from the requester's own claims, not client input.
type CreateInput struct {
	SedeID            int64
	NombreComercial   string
	NombreAlternativo *string
	Genero            string
	GramosMinimo      string
	NumeroGenero      int32
	RequesterID       int64
	IP                string
	UserAgent         string
}

// Create registers a new fragancia and seeds its stock (vitrina + bodega,
// both zero) atomically.
func (s *Service) Create(ctx context.Context, in CreateInput) (generated.GetFraganciaByIDRow, error) {
	exists, err := s.Fragancias.ExistsNombreComercial(ctx, in.SedeID, in.NombreComercial, 0)
	if err != nil {
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}
	if exists {
		return generated.GetFraganciaByIDRow{}, domainerrors.NewConflict(
			"Nombre en uso",
			"Ya existe una fragancia con ese nombre comercial en esta sede.",
		)
	}

	numeroExists, err := s.Fragancias.ExistsNumeroGenero(ctx, in.SedeID, in.Genero, in.NumeroGenero, 0)
	if err != nil {
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}
	if numeroExists {
		return generated.GetFraganciaByIDRow{}, domainerrors.NewConflict(
			"Número en uso",
			"Ya existe una fragancia con ese número en este género.",
		)
	}

	var fragancia generated.Fragancia
	err = commonrepo.WithTx(ctx, s.Pool, func(tx pgx.Tx) error {
		txFragancias := repo.NewPostgres(tx)
		txStock := stockactual.NewPostgres(tx)

		f, err := txFragancias.Insert(ctx, in.SedeID, in.NombreComercial, in.NombreAlternativo, in.Genero, in.GramosMinimo, in.NumeroGenero)
		if err != nil {
			return err
		}
		fragancia = f

		return txStock.InitializeStock(ctx, in.SedeID, stockactual.TipoItemFragancia, f.ID)
	})
	if err != nil {
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}

	result, err := s.Fragancias.GetByID(ctx, fragancia.ID)
	if err != nil {
		return generated.GetFraganciaByIDRow{}, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "fragancia_creada", in.IP, in.UserAgent, &fragancia.ID, nil, snapshot(fragancia))

	return result, nil
}
