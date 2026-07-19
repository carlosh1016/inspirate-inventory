package movimientos

import (
	"context"
	"strings"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
)

// DanadoInput is the request payload plus the requester's context. A single
// item per request (not a batch, unlike entrada/traslado).
type DanadoInput struct {
	SedeID      int64
	TipoItem    string
	ItemID      int64
	Ubicacion   string
	Cantidad    string
	Motivo      string
	RequesterID int64
	IP          string
	UserAgent   string
}

// Danado reduces stock for one item, requiring a motivo (also enforced by
// the DB's own CHECK constraint on movimientos_inventario). Not audited
// (operational action, not a deliberate inventory correction).
func (s *Service) Danado(ctx context.Context, in DanadoInput) (domainmovimientos.Movimiento, error) {
	cantidad, err := parsePositiveDecimal(in.Cantidad)
	if err != nil {
		return domainmovimientos.Movimiento{}, err
	}

	motivo := strings.TrimSpace(in.Motivo)
	if motivo == "" {
		return domainmovimientos.Movimiento{}, domainerrors.NewValidation(
			"Motivo requerido",
			"Debes indicar un motivo para registrar mercancía dañada.",
			nil,
		)
	}

	result, err := s.RegisterBatch(ctx, in.SedeID, in.RequesterID, []domainmovimientos.MovimientoInput{
		{
			TipoItem:  domainstock.TipoItem(in.TipoItem),
			ItemID:    in.ItemID,
			Tipo:      domainmovimientos.TipoDanado,
			Ubicacion: domainstock.Ubicacion(in.Ubicacion),
			Cantidad:  cantidad.Neg(),
			Motivo:    &motivo,
		},
	})
	if err != nil {
		return domainmovimientos.Movimiento{}, err
	}

	return result[0], nil
}
