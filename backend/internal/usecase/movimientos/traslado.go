package movimientos

import (
	"context"

	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
)

// TrasladoItem is one item in a POST /movimientos/traslado batch. Ubicación
// is always bodega -> vitrina, so it isn't part of the request shape.
type TrasladoItem struct {
	TipoItem string
	ItemID   int64
	Cantidad string
}

// TrasladoInput is the request payload plus the requester's context.
type TrasladoInput struct {
	SedeID      int64
	Items       []TrasladoItem
	RequesterID int64
	IP          string
	UserAgent   string
}

// Traslado moves stock from bodega to vitrina for each item in the batch,
// generating two movimientos per item (a bodega salida and a vitrina
// entrada). Both go through the same RegisterBatch transaction, so
// Postgres's per-transaction NOW() gives them an identical created_at —
// there's no need to compute or pass a timestamp explicitly. Not audited
// (high-volume operational action).
func (s *Service) Traslado(ctx context.Context, in TrasladoInput) ([]domainmovimientos.Movimiento, error) {
	inputs := make([]domainmovimientos.MovimientoInput, 0, len(in.Items)*2)
	for _, item := range in.Items {
		cantidad, err := parsePositiveDecimal(item.Cantidad)
		if err != nil {
			return nil, err
		}

		tipoItem := domainstock.TipoItem(item.TipoItem)
		inputs = append(inputs,
			domainmovimientos.MovimientoInput{
				TipoItem:  tipoItem,
				ItemID:    item.ItemID,
				Tipo:      domainmovimientos.TipoTrasladoBodegaVitrina,
				Ubicacion: domainstock.UbicacionBodega,
				Cantidad:  cantidad.Neg(),
			},
			domainmovimientos.MovimientoInput{
				TipoItem:  tipoItem,
				ItemID:    item.ItemID,
				Tipo:      domainmovimientos.TipoTrasladoBodegaVitrina,
				Ubicacion: domainstock.UbicacionVitrina,
				Cantidad:  cantidad,
			},
		)
	}

	return s.RegisterBatch(ctx, in.SedeID, in.RequesterID, inputs)
}
