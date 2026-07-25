// Package reportes is the persistence port for report data. It exposes one
// read method per report dataset; the port's result types are plain data
// structs (following the auditoria.Entry precedent), consumed by
// usecase/reportes to build the XLSX files.
package reportes

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// RangoFiltro scopes a time-ranged report to a sede and [Desde, Hasta] window,
// with an optional usuario filter (UsuarioID == 0 means "todos").
type RangoFiltro struct {
	SedeID    int64
	Desde     time.Time
	Hasta     time.Time
	UsuarioID int64
}

// StockFiltro scopes the stock snapshot report (no date range).
type StockFiltro struct {
	SedeID           int64
	IncluirInactivos bool
}

// MovimientosFiltro scopes the movimientos report.
type MovimientosFiltro struct {
	SedeID    int64
	Desde     time.Time
	Hasta     time.Time
	Tipo      string // "" = todos
	TipoItem  string // "" = todos
	ItemID    int64  // 0 = todos
	UsuarioID int64  // 0 = todos
}

// VentasResumen holds the aggregate totals for the ventas report summary sheet.
type VentasResumen struct {
	TotalEfectivo      decimal.Decimal
	TotalNequi         decimal.Decimal
	TotalDaviplata     decimal.Decimal
	TotalTransferencia decimal.Decimal
	TotalOtros         decimal.Decimal
	VentasCount        int64
	TotalVentas        decimal.Decimal
	DescuentoTotal     decimal.Decimal
}

// VentasPorVendedora is one vendedora's ventas count and total.
type VentasPorVendedora struct {
	UsuarioID      int64
	NombreCompleto string
	VentasCount    int64
	Total          decimal.Decimal
}

// VentaDetalle is one row of the ventas detail sheet.
type VentaDetalle struct {
	ID               int64
	CreatedAt        time.Time
	UsuarioNombre    string
	MetodoPagoNombre string
	Subtotal         decimal.Decimal
	DescuentoPct     decimal.Decimal
	DescuentoMonto   decimal.Decimal
	Total            decimal.Decimal
	Observaciones    *string
}

// VentaItem is one row of the ventas items sheet.
type VentaItem struct {
	VentaID         int64
	CreatedAt       time.Time
	TipoLinea       string
	FraganciaNombre *string
	EnvaseNombre    string
	ProductoNombre  *string
	FeromonaNombre  *string
	Gramos          *decimal.Decimal
	Cantidad        int32
	PrecioUnitario  decimal.Decimal
	Subtotal        decimal.Decimal
}

// StockFragancia is one row of the stock fragancias sheet.
type StockFragancia struct {
	NombreComercial   string
	NombreAlternativo *string
	Genero            string
	StockVitrina      decimal.Decimal
	StockBodega       decimal.Decimal
	StockTotal        decimal.Decimal
	Minimo            decimal.Decimal
	BajoMinimo        bool
}

// StockEnvase is one row of the stock envases sheet.
type StockEnvase struct {
	Tipo               string
	TamanoOz           decimal.Decimal
	Color              string
	PrecioSolo         decimal.Decimal
	PrecioConFragancia decimal.Decimal
	PrecioRecarga      decimal.Decimal
	StockVitrina       decimal.Decimal
	StockBodega        decimal.Decimal
	StockTotal         decimal.Decimal
	Minimo             decimal.Decimal
	BajoMinimo         bool
}

// StockProducto is one row of the stock productos sheet.
type StockProducto struct {
	Nombre       string
	Categoria    string
	Precio       decimal.Decimal
	StockVitrina decimal.Decimal
	StockBodega  decimal.Decimal
	StockTotal   decimal.Decimal
	Minimo       decimal.Decimal
	BajoMinimo   bool
}

// StockAlerta is one row of the stock alertas sheet (items bajo mínimo).
type StockAlerta struct {
	Tipo        string
	Nombre      string
	StockActual decimal.Decimal
	Minimo      decimal.Decimal
	Faltante    decimal.Decimal
}

// Movimiento is one row of the movimientos sheet.
type Movimiento struct {
	CreatedAt      time.Time
	Tipo           string
	TipoItem       string
	ItemNombre     string
	Ubicacion      string
	Cantidad       decimal.Decimal
	StockAnterior  decimal.Decimal
	StockPosterior decimal.Decimal
	UsuarioNombre  string
	Motivo         *string
	VentaID        *int64
}

// CuadreCerrado is one row of the cuadres sheet (only cerrados).
type CuadreCerrado struct {
	Fecha               time.Time
	Estado              string
	FondoBase           decimal.Decimal
	TotalEfectivo       decimal.Decimal
	TotalNequi          decimal.Decimal
	TotalDaviplata      decimal.Decimal
	TotalTransferencia  decimal.Decimal
	TotalOtros          decimal.Decimal
	TotalPagos          decimal.Decimal
	TotalConsignaciones decimal.Decimal
	ValorTurno          decimal.Decimal
	SaldoCalculado      decimal.Decimal
	CerradoPor          *string
	CerradoAt           *time.Time
	Observaciones       *string
}

// CuadrePago is one row of the cuadres pagos sheet.
type CuadrePago struct {
	FechaCuadre   time.Time
	Concepto      string
	Monto         decimal.Decimal
	UsuarioNombre string
	CreatedAt     time.Time
}

// CuadreConsignacion is one row of the cuadres consignaciones sheet.
type CuadreConsignacion struct {
	FechaCuadre   time.Time
	Monto         decimal.Decimal
	Banco         *string
	Referencia    *string
	UsuarioNombre string
	CreatedAt     time.Time
}

// SesionResumen is one vendedora's aggregate in the sesiones summary sheet.
type SesionResumen struct {
	UsuarioID      int64
	NombreCompleto string
	TotalHoras     *time.Duration
	DiasTrabajados int64
	SesionesCount  int64
}

// SesionDetalle is one row of the sesiones detail sheet.
type SesionDetalle struct {
	NombreCompleto  string
	EntradaAt       time.Time
	SalidaAt        *time.Time
	HorasTrabajadas *time.Duration
}

// Repository is the persistence port for report datasets.
type Repository interface {
	VentasResumen(ctx context.Context, f RangoFiltro) (VentasResumen, error)
	VentasPorVendedora(ctx context.Context, f RangoFiltro) ([]VentasPorVendedora, error)
	VentasDetalle(ctx context.Context, f RangoFiltro) ([]VentaDetalle, error)
	VentasItems(ctx context.Context, f RangoFiltro) ([]VentaItem, error)

	StockFragancias(ctx context.Context, f StockFiltro) ([]StockFragancia, error)
	StockEnvases(ctx context.Context, f StockFiltro) ([]StockEnvase, error)
	StockProductos(ctx context.Context, f StockFiltro) ([]StockProducto, error)
	StockAlertas(ctx context.Context, f StockFiltro) ([]StockAlerta, error)

	Movimientos(ctx context.Context, f MovimientosFiltro) ([]Movimiento, error)

	CuadresCerrados(ctx context.Context, f RangoFiltro) ([]CuadreCerrado, error)
	CuadresPagos(ctx context.Context, f RangoFiltro) ([]CuadrePago, error)
	CuadresConsignaciones(ctx context.Context, f RangoFiltro) ([]CuadreConsignacion, error)

	SesionesResumen(ctx context.Context, f RangoFiltro) ([]SesionResumen, error)
	SesionesDetalle(ctx context.Context, f RangoFiltro) ([]SesionDetalle, error)
}
