package conteos

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	domainconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/conteos"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
)

// CrearInput is the request payload plus the requester's context. Periodo,
// when nil, defaults to the current month.
type CrearInput struct {
	SedeID      int64
	Periodo     *time.Time
	RequesterID int64
	IP          string
	UserAgent   string
}

type auditCrearSnapshot struct {
	Periodo string `json:"periodo"`
}

// Crear opens a new conteo mensual for SedeID's periodo (the first day of
// the target month), generating one row per active fragancia with its
// saldo_inicial/gramos_sistema already resolved. Only one conteo per (sede,
// periodo) may ever exist (enforced by the UNIQUE constraint too) — a
// conteo already existing for that month is a 409.
func (s *Service) Crear(ctx context.Context, in CrearInput) (*domainconteos.ConteoInventario, error) {
	periodo := s.primerDiaMesActual()
	if in.Periodo != nil {
		p := *in.Periodo
		periodo = time.Date(p.Year(), p.Month(), 1, 0, 0, 0, 0, s.Location)
	}

	if _, err := s.Conteos.GetBySedePeriodo(ctx, in.SedeID, periodo); err == nil {
		return nil, domainerrors.NewConflict("Ya existe un conteo para este mes", "Ya se creó un conteo de inventario para este periodo.")
	} else if !errors.Is(err, conteosrepo.ErrNotFound) {
		return nil, internalErr(err)
	}

	var conteoID int64
	err := commonrepo.WithTx(ctx, s.Pool, func(tx pgx.Tx) error {
		txConteos := conteosrepo.NewPostgres(tx)
		header, err := txConteos.Insert(ctx, in.SedeID, periodo, in.RequesterID)
		if err != nil {
			return err
		}
		conteoID = header.ID
		return txConteos.GenerarItems(ctx, header.ID, in.SedeID, periodo)
	})
	if err != nil {
		return nil, internalErr(err)
	}

	s.audit(ctx, &in.RequesterID, "conteo_creado", in.IP, in.UserAgent, &conteoID, nil, auditCrearSnapshot{
		Periodo: periodo.Format("2006-01-02"),
	})

	return s.GetByID(ctx, conteoID)
}
