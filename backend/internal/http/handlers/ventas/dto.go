package ventas

import (
	"fmt"
	"time"

	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// VentaUsuarioBrief is the (id, nombre_completo) of who registered a venta.
type VentaUsuarioBrief struct {
	ID             int64  `json:"id"`
	NombreCompleto string `json:"nombre_completo"`
}

// VentaMetodoPagoBrief is the (id, nombre, codigo) of a venta's método de
// pago.
type VentaMetodoPagoBrief struct {
	ID     int64  `json:"id"`
	Nombre string `json:"nombre"`
	Codigo string `json:"codigo"`
}

// VentaItemDetalleResponse is one line of a venta, with catalog references
// resolved to display names so the frontend doesn't have to.
type VentaItemDetalleResponse struct {
	ID                   int64     `json:"id"`
	TipoLinea            string    `json:"tipo_linea"`
	FraganciaID          *int64    `json:"fragancia_id,omitempty"`
	FraganciaNombre      *string   `json:"fragancia_nombre,omitempty"`
	VarianteEnvaseID     *int64    `json:"variante_envase_id,omitempty"`
	VarianteEnvaseNombre *string   `json:"variante_envase_nombre,omitempty"`
	ProductoID           *int64    `json:"producto_id,omitempty"`
	ProductoNombre       *string   `json:"producto_nombre,omitempty"`
	FeromonaProductoID   *int64    `json:"feromona_producto_id,omitempty"`
	FeromonaNombre       *string   `json:"feromona_nombre,omitempty"`
	GramosFragancia      *string   `json:"gramos_fragancia,omitempty"`
	Cantidad             int32     `json:"cantidad"`
	PrecioUnitario       string    `json:"precio_unitario"`
	Subtotal             string    `json:"subtotal"`
	CreatedAt            time.Time `json:"created_at"`
}

// VentaDetalladaResponse is the full response shape for POST/GET/PATCH
// /ventas/:id — Consecutivo is "V-" + id padded to 6 digits.
type VentaDetalladaResponse struct {
	ID                   int64                      `json:"id"`
	Consecutivo          string                     `json:"consecutivo"`
	Usuario              VentaUsuarioBrief          `json:"usuario"`
	MetodoPago           VentaMetodoPagoBrief       `json:"metodo_pago"`
	Items                []VentaItemDetalleResponse `json:"items"`
	Subtotal             string                     `json:"subtotal"`
	DescuentoPct         string                     `json:"descuento_pct"`
	DescuentoMonto       string                     `json:"descuento_monto"`
	Total                string                     `json:"total"`
	Observaciones        *string                    `json:"observaciones"`
	CreatedAt            time.Time                  `json:"created_at"`
	MovimientosGenerados []int64                    `json:"movimientos_generados"`
}

func formatConsecutivo(id int64) string {
	return fmt.Sprintf("V-%06d", id)
}

func toVentaItemDetalleResponse(item domainventas.VentaItem) VentaItemDetalleResponse {
	var gramos *string
	if item.GramosFragancia != nil {
		s := item.GramosFragancia.String()
		gramos = &s
	}
	return VentaItemDetalleResponse{
		ID:                   item.ID,
		TipoLinea:            string(item.TipoLinea),
		FraganciaID:          item.FraganciaID,
		FraganciaNombre:      item.FraganciaNombre,
		VarianteEnvaseID:     item.VarianteEnvaseID,
		VarianteEnvaseNombre: item.VarianteEnvaseNombre,
		ProductoID:           item.ProductoID,
		ProductoNombre:       item.ProductoNombre,
		FeromonaProductoID:   item.FeromonaProductoID,
		FeromonaNombre:       item.FeromonaNombre,
		GramosFragancia:      gramos,
		Cantidad:             item.Cantidad,
		PrecioUnitario:       item.PrecioUnitario.String(),
		Subtotal:             item.Subtotal.String(),
		CreatedAt:            item.CreatedAt,
	}
}

// toVentaDetalladaResponse builds the full response for a venta. Used for
// POST (fresh creation), GET (single venta, with movimientosGenerados
// filled in), and PATCH (movimientosGenerados left empty — the caller
// doesn't need it there, but the field stays present for a consistent
// shape).
func toVentaDetalladaResponse(v domainventas.Venta, movimientosGenerados []int64) VentaDetalladaResponse {
	items := make([]VentaItemDetalleResponse, len(v.Items))
	for i, item := range v.Items {
		items[i] = toVentaItemDetalleResponse(item)
	}
	if movimientosGenerados == nil {
		movimientosGenerados = []int64{}
	}
	return VentaDetalladaResponse{
		ID:                   v.ID,
		Consecutivo:          formatConsecutivo(v.ID),
		Usuario:              VentaUsuarioBrief{ID: v.UsuarioID, NombreCompleto: v.UsuarioNombre},
		MetodoPago:           VentaMetodoPagoBrief{ID: v.MetodoPagoID, Nombre: v.MetodoPagoNombre, Codigo: v.MetodoPagoCodigo},
		Items:                items,
		Subtotal:             v.Subtotal.String(),
		DescuentoPct:         v.DescuentoPct.String(),
		DescuentoMonto:       v.DescuentoMonto.String(),
		Total:                v.Total.String(),
		Observaciones:        v.Observaciones,
		CreatedAt:            v.CreatedAt,
		MovimientosGenerados: movimientosGenerados,
	}
}

// VentaListItemResponse is one row of GET /ventas — lighter than the
// detailed response, no items array (fetch GET /ventas/:id for that).
type VentaListItemResponse struct {
	ID               int64     `json:"id"`
	Consecutivo      string    `json:"consecutivo"`
	UsuarioID        int64     `json:"usuario_id"`
	UsuarioNombre    string    `json:"usuario_nombre"`
	MetodoPagoID     int64     `json:"metodo_pago_id"`
	MetodoPagoNombre string    `json:"metodo_pago_nombre"`
	MetodoPagoCodigo string    `json:"metodo_pago_codigo"`
	ItemsCount       int64     `json:"items_count"`
	Subtotal         string    `json:"subtotal"`
	DescuentoPct     string    `json:"descuento_pct"`
	DescuentoMonto   string    `json:"descuento_monto"`
	Total            string    `json:"total"`
	Observaciones    *string   `json:"observaciones,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

func toVentaListItemResponse(row generated.ListVentasPaginatedRow) VentaListItemResponse {
	var observaciones *string
	if row.Observaciones.Valid {
		observaciones = &row.Observaciones.String
	}
	return VentaListItemResponse{
		ID:               row.ID,
		Consecutivo:      formatConsecutivo(row.ID),
		UsuarioID:        row.UsuarioID,
		UsuarioNombre:    row.UsuarioNombre,
		MetodoPagoID:     row.MetodoPagoID,
		MetodoPagoNombre: row.MetodoPagoNombre,
		MetodoPagoCodigo: row.MetodoPagoCodigo,
		ItemsCount:       row.ItemsCount,
		Subtotal:         row.Subtotal.String(),
		DescuentoPct:     row.DescuentoPct.String(),
		DescuentoMonto:   row.DescuentoMonto.String(),
		Total:            row.Total.String(),
		Observaciones:    observaciones,
		CreatedAt:        row.CreatedAt.Time,
	}
}

// ResumenPorMetodoPagoResponse breaks today's total down by método de pago.
type ResumenPorMetodoPagoResponse struct {
	Efectivo      string `json:"efectivo"`
	Nequi         string `json:"nequi"`
	Daviplata     string `json:"daviplata"`
	Transferencia string `json:"transferencia"`
	Otros         string `json:"otros"`
}

// ResumenVendedoraResponse is one vendedora's share of today's ventas.
type ResumenVendedoraResponse struct {
	UsuarioID      int64  `json:"usuario_id"`
	NombreCompleto string `json:"nombre_completo"`
	VentasCount    int64  `json:"ventas_count"`
	Total          string `json:"total"`
}

// ResumenFraganciaResponse is one fragancia's share of today's sales.
type ResumenFraganciaResponse struct {
	ID              int64  `json:"id"`
	NombreComercial string `json:"nombre_comercial"`
	GramosVendidos  string `json:"gramos_vendidos"`
	MontoVendido    string `json:"monto_vendido"`
}

// ResumenHoyResponse is the full response for GET /ventas/hoy/resumen. The
// exact shape isn't dictated by the spec (it references a DTO from an
// unrelated draft) — this is derived directly from the 3 queries the spec
// did provide.
type ResumenHoyResponse struct {
	Fecha          string                       `json:"fecha"`
	VentasCount    int64                        `json:"ventas_count"`
	TotalDia       string                       `json:"total_dia"`
	DescuentoTotal string                       `json:"descuento_total"`
	PorMetodoPago  ResumenPorMetodoPagoResponse `json:"por_metodo_pago"`
	PorVendedora   []ResumenVendedoraResponse   `json:"por_vendedora"`
	TopFragancias  []ResumenFraganciaResponse   `json:"top_fragancias"`
}

// CreateVentaItemRequest is one line in a POST /ventas request.
type CreateVentaItemRequest struct {
	TipoLinea          string  `json:"tipo_linea" validate:"required,oneof=envase_con_fragancia recarga envase_solo producto_otro"`
	FraganciaID        *int64  `json:"fragancia_id,omitempty" validate:"omitempty,gt=0"`
	VarianteEnvaseID   *int64  `json:"variante_envase_id,omitempty" validate:"omitempty,gt=0"`
	ProductoID         *int64  `json:"producto_id,omitempty" validate:"omitempty,gt=0"`
	FeromonaProductoID *int64  `json:"feromona_producto_id,omitempty" validate:"omitempty,gt=0"`
	GramosFragancia    *string `json:"gramos_fragancia,omitempty" validate:"omitempty,numeric"`
	Cantidad           int32   `json:"cantidad" validate:"required,gt=0"`
}

// CreateVentaRequest is the payload for POST /ventas. Deep, per-tipo_linea
// coherence isn't checked by struct tags (dive only validates shallow
// per-field tags without positional information) — the usecase does that by
// hand and reports failures with an index in DomainError.Extra.
type CreateVentaRequest struct {
	MetodoPagoID  int64                    `json:"metodo_pago_id" validate:"required,gt=0"`
	Observaciones *string                  `json:"observaciones,omitempty" validate:"omitempty,max=1000"`
	Items         []CreateVentaItemRequest `json:"items" validate:"required,min=1,max=50,dive"`
}

// UpdateVentaRequest is the payload for PATCH /ventas/:id. The handler
// decodes with DisallowUnknownFields, so any field beyond Observaciones
// (e.g. {"total": 0}) fails as an invalid body rather than being silently
// ignored.
type UpdateVentaRequest struct {
	Observaciones *string `json:"observaciones" validate:"omitempty,max=1000"`
}
