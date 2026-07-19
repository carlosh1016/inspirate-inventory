package movimientos

import (
	"context"

	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
)

// EntradaMercanciaItem is one item in a POST /movimientos/entrada-mercancia
// batch. Cantidad is a positive decimal string.
type EntradaMercanciaItem struct {
	TipoItem  string
	ItemID    int64
	Ubicacion string
	Cantidad  string
	Motivo    *string
}

// EntradaMercanciaInput is the request payload plus the requester's
// context. SedeID is inherited from the requester's own claims.
type EntradaMercanciaInput struct {
	SedeID      int64
	Items       []EntradaMercanciaItem
	RequesterID int64
	IP          string
	UserAgent   string
}

// EntradaMercancia registers new stock arriving into the given ubicación
// for each item in the batch. Not audited (high-volume operational action).
func (s *Service) EntradaMercancia(ctx context.Context, in EntradaMercanciaInput) ([]domainmovimientos.Movimiento, error) {
	inputs := make([]domainmovimientos.MovimientoInput, len(in.Items))
	for i, item := range in.Items {
		cantidad, err := parsePositiveDecimal(item.Cantidad)
		if err != nil {
			return nil, err
		}
		inputs[i] = domainmovimientos.MovimientoInput{
			TipoItem:  domainstock.TipoItem(item.TipoItem),
			ItemID:    item.ItemID,
			Tipo:      domainmovimientos.TipoEntradaMercancia,
			Ubicacion: domainstock.Ubicacion(item.Ubicacion),
			Cantidad:  cantidad,
			Motivo:    item.Motivo,
		}
	}

	return s.RegisterBatch(ctx, in.SedeID, in.RequesterID, inputs)
}
