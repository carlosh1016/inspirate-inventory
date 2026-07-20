// Package ventas contains pure domain logic for ventas: no I/O, no HTTP, no
// SQL. Pricing and discount are pure calculators here too — the
// transactional orchestration that loads catalog entities, consolidates
// stock movements, and persists everything lives in usecase/ventas.
package ventas

import (
	"time"

	"github.com/shopspring/decimal"
)

// TipoLinea identifies the kind of venta_item, matching the tipo_linea_enum
// in the schema. Each value implies a fixed set of non-null FK columns,
// enforced both by chk_venta_items_tipo_linea in the DB and by usecase-level
// validation before any DB write.
type TipoLinea string

const (
	TipoLineaEnvaseConFragancia TipoLinea = "envase_con_fragancia"
	TipoLineaRecarga            TipoLinea = "recarga"
	TipoLineaEnvaseSolo         TipoLinea = "envase_solo"
	TipoLineaProductoOtro       TipoLinea = "producto_otro"
)

// Venta is the domain representation of an immutable sale (except
// Observaciones, which admin may edit post-hoc). Subtotal/DescuentoPct/
// DescuentoMonto/Total are pre-calculated at creation time and never
// recomputed — historical ventas don't change if discount rules change
// tomorrow. UsuarioNombre/MetodoPagoNombre/MetodoPagoCodigo are denormalized
// display fields that come along whenever a Venta is loaded via GetByID —
// not written, never used for business logic, only for the HTTP response.
type Venta struct {
	ID               int64
	SedeID           int64
	UsuarioID        int64
	UsuarioNombre    string
	MetodoPagoID     int64
	MetodoPagoNombre string
	MetodoPagoCodigo string
	Subtotal         decimal.Decimal
	DescuentoPct     decimal.Decimal
	DescuentoMonto   decimal.Decimal
	Total            decimal.Decimal
	Observaciones    *string
	CreatedAt        time.Time
	Items            []VentaItem
}

// VentaItem is one line of a venta. Which of FraganciaID/VarianteEnvaseID/
// ProductoID/FeromonaProductoID/GramosFragancia are non-nil is dictated by
// TipoLinea (see chk_venta_items_tipo_linea). The *Nombre fields are
// denormalized display data, same caveat as Venta's.
type VentaItem struct {
	ID                   int64
	VentaID              int64
	TipoLinea            TipoLinea
	FraganciaID          *int64
	FraganciaNombre      *string
	VarianteEnvaseID     *int64
	VarianteEnvaseNombre *string
	ProductoID           *int64
	ProductoNombre       *string
	FeromonaProductoID   *int64
	FeromonaNombre       *string
	GramosFragancia      *decimal.Decimal
	Cantidad             int32
	PrecioUnitario       decimal.Decimal
	Subtotal             decimal.Decimal
	CreatedAt            time.Time
}
