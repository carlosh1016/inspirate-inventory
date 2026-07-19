package movimientos

import (
	"time"

	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	commonrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// MovimientoItemResponse is the (id, nombre) of the item a movimiento
// affected.
type MovimientoItemResponse struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
}

// MovimientoUsuarioResponse is the (id, nombre_completo) of the user who
// registered a movimiento. NombreCompleto is only populated for the List
// endpoint (which already joins usuarios) — POST responses only have the
// id on hand, since the acting user already knows their own name and
// fetching it just for the echo isn't worth an extra query.
type MovimientoUsuarioResponse struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo,omitempty"`
}

// MovimientoResponse is one row of movimientos_inventario.
type MovimientoResponse struct {
	ID             int64                     `json:"id"`
	TipoItem       string                    `json:"tipo_item"`
	Item           MovimientoItemResponse    `json:"item"`
	Tipo           string                    `json:"tipo"`
	Ubicacion      string                    `json:"ubicacion"`
	Cantidad       string                    `json:"cantidad"`
	StockAnterior  string                    `json:"stock_anterior"`
	StockPosterior string                    `json:"stock_posterior"`
	Usuario        MovimientoUsuarioResponse `json:"usuario"`
	VentaID        *int64                    `json:"venta_id,omitempty"`
	Motivo         *string                   `json:"motivo,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
}

func toMovimientoResponseFromDomain(m domainmovimientos.Movimiento) MovimientoResponse {
	return MovimientoResponse{
		ID:             m.ID,
		TipoItem:       string(m.TipoItem),
		Item:           MovimientoItemResponse{ID: m.ItemID, Nombre: m.ItemNombre},
		Tipo:           string(m.Tipo),
		Ubicacion:      string(m.Ubicacion),
		Cantidad:       m.Cantidad.String(),
		StockAnterior:  m.StockAnterior.String(),
		StockPosterior: m.StockPosterior.String(),
		Usuario:        MovimientoUsuarioResponse{ID: m.UsuarioID},
		VentaID:        m.VentaID,
		Motivo:         m.Motivo,
		CreatedAt:      m.CreatedAt,
	}
}

func toMovimientoResponseFromList(row generated.ListMovimientosPaginatedRow) MovimientoResponse {
	return MovimientoResponse{
		ID:             row.ID,
		TipoItem:       string(row.TipoItem),
		Item:           MovimientoItemResponse{ID: row.ItemID, Nombre: row.ItemNombre},
		Tipo:           string(row.Tipo),
		Ubicacion:      string(row.Ubicacion),
		Cantidad:       row.Cantidad.String(),
		StockAnterior:  row.StockAnterior.String(),
		StockPosterior: row.StockPosterior.String(),
		Usuario:        MovimientoUsuarioResponse{ID: row.UsuarioID, NombreCompleto: row.UsuarioNombre},
		VentaID:        commonrepo.Int8Ptr(row.VentaID),
		Motivo:         commonrepo.StringPtr(row.Motivo),
		CreatedAt:      row.CreatedAt.Time,
	}
}

// EntradaMercanciaItemRequest is one item in a POST
// /movimientos/entrada-mercancia batch.
type EntradaMercanciaItemRequest struct {
	TipoItem  string  `json:"tipo_item" validate:"required,oneof=fragancia variante_envase producto"`
	ItemID    int64   `json:"item_id" validate:"required,gt=0"`
	Ubicacion string  `json:"ubicacion" validate:"required,oneof=vitrina bodega"`
	Cantidad  string  `json:"cantidad" validate:"required,numeric"`
	Motivo    *string `json:"motivo,omitempty" validate:"omitempty,max=500"`
}

// EntradaMercanciaRequest is the payload for POST
// /movimientos/entrada-mercancia.
type EntradaMercanciaRequest struct {
	Items []EntradaMercanciaItemRequest `json:"items" validate:"required,min=1,max=100,dive"`
}

// TrasladoItemRequest is one item in a POST /movimientos/traslado batch.
// Ubicación isn't part of the request — traslado always moves
// bodega -> vitrina.
type TrasladoItemRequest struct {
	TipoItem string `json:"tipo_item" validate:"required,oneof=fragancia variante_envase producto"`
	ItemID   int64  `json:"item_id" validate:"required,gt=0"`
	Cantidad string `json:"cantidad" validate:"required,numeric"`
}

// TrasladoRequest is the payload for POST /movimientos/traslado.
type TrasladoRequest struct {
	Items []TrasladoItemRequest `json:"items" validate:"required,min=1,max=100,dive"`
}

// DanadoRequest is the payload for POST /movimientos/danado.
type DanadoRequest struct {
	TipoItem  string `json:"tipo_item" validate:"required,oneof=fragancia variante_envase producto"`
	ItemID    int64  `json:"item_id" validate:"required,gt=0"`
	Ubicacion string `json:"ubicacion" validate:"required,oneof=vitrina bodega"`
	Cantidad  string `json:"cantidad" validate:"required,numeric"`
	Motivo    string `json:"motivo" validate:"required,min=1,max=500"`
}

// AjusteRequest is the payload for POST /movimientos/ajuste.
// CantidadNueva is the absolute target quantity, not a delta.
type AjusteRequest struct {
	TipoItem      string `json:"tipo_item" validate:"required,oneof=fragancia variante_envase producto"`
	ItemID        int64  `json:"item_id" validate:"required,gt=0"`
	Ubicacion     string `json:"ubicacion" validate:"required,oneof=vitrina bodega"`
	CantidadNueva string `json:"cantidad_nueva" validate:"required,numeric"`
	Motivo        string `json:"motivo" validate:"required,min=1,max=500"`
}

// CorreccionRequest has the same shape as AjusteRequest.
type CorreccionRequest = AjusteRequest
