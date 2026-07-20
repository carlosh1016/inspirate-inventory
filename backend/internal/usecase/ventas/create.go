package ventas

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/envases"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainmovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/movimientos"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/domain/productos"
	domainstock "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/stock"
	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	fraganciasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
	metodospagorepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
	modelosenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	productosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	variantesenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	ventaitemsrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/venta_items"
	ventasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/ventas"
)

const (
	minVentaItems = 1
	maxVentaItems = 50
)

// CreateVentaItemInput is one line of a POST /ventas request. Which pointer
// fields are set must match TipoLinea (see validateItemCoherence) —
// mirrors chk_venta_items_tipo_linea in the DB.
type CreateVentaItemInput struct {
	TipoLinea          domainventas.TipoLinea
	FraganciaID        *int64
	VarianteEnvaseID   *int64
	ProductoID         *int64
	FeromonaProductoID *int64
	GramosFragancia    *decimal.Decimal
	Cantidad           int32
}

// CreateVentaInput is the request payload plus the requester's context.
// IdempotencyKey is nil when the header was absent or empty (treated the
// same). RequestBodyHash is only meaningful when IdempotencyKey != nil.
type CreateVentaInput struct {
	SedeID          int64
	UsuarioID       int64
	MetodoPagoID    int64
	Items           []CreateVentaItemInput
	Observaciones   *string
	IdempotencyKey  *string
	RequestBodyHash string
	IP              string
	UserAgent       string
}

// CreateVentaOutput is either a freshly created venta, or — when
// FromIdempotency is true — the verbatim cached response (CachedStatus,
// CachedBody) from an earlier identical request, to be replayed as-is.
type CreateVentaOutput struct {
	Venta           *domainventas.Venta
	FromIdempotency bool
	CachedStatus    int32
	CachedBody      []byte
}

// CreateVenta is the core usecase of the whole system: it prices every
// line, applies the automatic discount, consolidates stock outflows across
// lines, and registers the venta + its movimientos atomically. See the
// package doc and the tanda plan for the full step-by-step.
func (s *Service) CreateVenta(ctx context.Context, in CreateVentaInput) (*CreateVentaOutput, error) {
	// Fase 1: pre-validación, fuera de transacción.
	if in.IdempotencyKey != nil {
		cached, err := s.checkIdempotency(ctx, *in.IdempotencyKey, in.UsuarioID, in.RequestBodyHash)
		if err != nil {
			return nil, err
		}
		if cached != nil {
			return &CreateVentaOutput{FromIdempotency: true, CachedStatus: cached.Status, CachedBody: cached.Body}, nil
		}
	}

	if err := s.validateMetodoPago(ctx, in.MetodoPagoID); err != nil {
		return nil, err
	}

	if len(in.Items) < minVentaItems || len(in.Items) > maxVentaItems {
		return nil, domainerrors.NewValidation("Solicitud inválida", "Debes incluir entre 1 y 50 ítems.", nil)
	}
	var coherenceErrors []domainventas.ItemErrorDetalle
	for i, item := range in.Items {
		if detalle := validateItemCoherence(item, i); detalle != nil {
			coherenceErrors = append(coherenceErrors, *detalle)
		}
	}
	if len(coherenceErrors) > 0 {
		return nil, itemErrorsToDomainError(coherenceErrors)
	}

	// Fase 2: dentro de una transacción RepeatableRead.
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, internalErr(err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // no-op once committed
	}()

	ventaID, err := s.createVentaTx(ctx, tx, in)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, internalErr(err)
	}

	venta, err := s.loadVentaCompleta(ctx, ventaID)
	if err != nil {
		return nil, err
	}

	return &CreateVentaOutput{Venta: venta, FromIdempotency: false}, nil
}

func (s *Service) validateMetodoPago(ctx context.Context, id int64) error {
	metodo, err := s.MetodosPago.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, metodospagorepo.ErrNotFound) {
			return domainerrors.NewValidation("Método de pago inválido", "El método de pago indicado no existe o está inactivo.", nil)
		}
		return internalErr(err)
	}
	if !metodo.Activo {
		return domainerrors.NewValidation("Método de pago inválido", "El método de pago indicado no existe o está inactivo.", nil)
	}
	return nil
}

// validateItemCoherence replicates chk_venta_items_tipo_linea in Go, before
// touching the DB at all.
func validateItemCoherence(item CreateVentaItemInput, index int) *domainventas.ItemErrorDetalle {
	motivo := ""
	switch item.TipoLinea {
	case domainventas.TipoLineaEnvaseConFragancia, domainventas.TipoLineaRecarga:
		motivo = validateFraganciaLine(item)
	case domainventas.TipoLineaEnvaseSolo:
		motivo = validateEnvaseSoloLine(item)
	case domainventas.TipoLineaProductoOtro:
		motivo = validateProductoOtroLine(item)
	default:
		motivo = "tipo_linea desconocido"
	}
	if motivo == "" && item.Cantidad <= 0 {
		motivo = "cantidad debe ser mayor que 0"
	}
	if motivo == "" {
		return nil
	}
	return &domainventas.ItemErrorDetalle{Index: index, Motivo: motivo}
}

func validateFraganciaLine(item CreateVentaItemInput) string {
	if item.FraganciaID == nil || item.VarianteEnvaseID == nil {
		return "fragancia_id y variante_envase_id son requeridos para este tipo de línea"
	}
	if item.ProductoID != nil {
		return "producto_id no debe enviarse para este tipo de línea"
	}
	if item.GramosFragancia == nil || !item.GramosFragancia.IsPositive() {
		return "gramos_fragancia debe ser mayor que 0"
	}
	return ""
}

func validateEnvaseSoloLine(item CreateVentaItemInput) string {
	if item.VarianteEnvaseID == nil {
		return "variante_envase_id es requerido para este tipo de línea"
	}
	if item.FraganciaID != nil || item.ProductoID != nil || item.GramosFragancia != nil || item.FeromonaProductoID != nil {
		return "este tipo de línea solo admite variante_envase_id"
	}
	return ""
}

func validateProductoOtroLine(item CreateVentaItemInput) string {
	if item.ProductoID == nil {
		return "producto_id es requerido para este tipo de línea"
	}
	if item.FraganciaID != nil || item.VarianteEnvaseID != nil || item.GramosFragancia != nil || item.FeromonaProductoID != nil {
		return "este tipo de línea solo admite producto_id"
	}
	return ""
}

func itemErrorsToDomainError(items []domainventas.ItemErrorDetalle) error {
	err := domainerrors.NewValidation(
		"Ítems inválidos",
		"Uno o más ítems de la venta no son válidos.",
		nil,
	)
	err.Extra = domainventas.ItemErrorExtra{Items: items}
	return err
}

// stockKey identifies one (tipo_item, item_id) pair for stock consolidation
// — every venta_item touching the same catalog item accumulates into a
// single MovimientoInput, keyed by this.
type stockKey struct {
	TipoItem string
	ItemID   int64
}

// lineResult is one priced, ready-to-insert venta_item.
type lineResult struct {
	item           CreateVentaItemInput
	precioUnitario decimal.Decimal
	subtotal       decimal.Decimal
}

// createVentaTx runs entirely inside tx: loads every referenced catalog
// entity, prices each line, applies the discount, consolidates stock
// outflows, inserts the venta + its items, registers the resulting
// movimientos via s.Movimientos.RegisterBatchTx, and — if IdempotencyKey is
// set — stores the response snapshot. Returns the new venta's ID.
func (s *Service) createVentaTx(ctx context.Context, tx pgx.Tx, in CreateVentaInput) (int64, error) {
	// TODO(cuadre): validar que el cuadre del día no esté cerrado (Tanda 5+).

	loader := newEntityLoader(tx)

	lines, subtotalVenta, movimientoInputs, loadErrors, err := s.processItems(ctx, loader, in.Items)
	if err != nil {
		return 0, err
	}
	if len(loadErrors) > 0 {
		return 0, itemErrorsToDomainError(loadErrors)
	}

	discount := s.Discount.Apply(subtotalVenta)

	venta, err := ventasrepo.NewPostgres(tx).Insert(ctx, ventasrepo.InsertParams{
		SedeID:         in.SedeID,
		UsuarioID:      in.UsuarioID,
		MetodoPagoID:   in.MetodoPagoID,
		Subtotal:       subtotalVenta,
		DescuentoPct:   discount.Pct,
		DescuentoMonto: discount.Monto,
		Total:          discount.Total,
		Observaciones:  in.Observaciones,
	})
	if err != nil {
		return 0, internalErr(err)
	}

	txVentaItems := ventaitemsrepo.NewPostgres(tx)
	for _, line := range lines {
		if _, err := txVentaItems.Insert(ctx, ventaitemsrepo.InsertParams{
			VentaID:            venta.ID,
			TipoLinea:          string(line.item.TipoLinea),
			FraganciaID:        line.item.FraganciaID,
			VarianteEnvaseID:   line.item.VarianteEnvaseID,
			ProductoID:         line.item.ProductoID,
			FeromonaProductoID: line.item.FeromonaProductoID,
			GramosFragancia:    line.item.GramosFragancia,
			Cantidad:           line.item.Cantidad,
			PrecioUnitario:     line.precioUnitario,
			Subtotal:           line.subtotal,
		}); err != nil {
			return 0, internalErr(err)
		}
	}

	for i := range movimientoInputs {
		movimientoInputs[i].VentaID = &venta.ID
	}
	if _, err := s.Movimientos.RegisterBatchTx(ctx, tx, in.SedeID, in.UsuarioID, movimientoInputs); err != nil {
		var domainErr *domainerrors.DomainError
		if errors.As(err, &domainErr) && domainErr.Code == domainerrors.CodeBusinessRule {
			if items, ok := domainErr.Extra.([]domainstock.StockInsuficienteItem); ok {
				return 0, stockInsuficienteToDomainError(items, indexByStockKey(in.Items))
			}
		}
		return 0, err
	}

	return venta.ID, nil
}

// stockDelta is one (tipo_item, item_id) quantity change contributed by a
// single line, before consolidation across the whole venta.
type stockDelta struct {
	TipoItem string
	ItemID   int64
	Delta    decimal.Decimal
}

// processItems loads every referenced catalog entity (memoized via loader
// — a fragancia/variante/producto referenced by several lines is only
// fetched once), prices each line, and consolidates stock outflows across
// all lines into one MovimientoInput per catalog item. Structural coherence
// was already validated in fase 1; this only checks that referenced
// entities exist, aren't soft-deleted, and are active.
func (s *Service) processItems(ctx context.Context, loader *entityLoader, items []CreateVentaItemInput) ([]lineResult, decimal.Decimal, []domainmovimientos.MovimientoInput, []domainventas.ItemErrorDetalle, error) {
	var lines []lineResult
	var loadErrors []domainventas.ItemErrorDetalle
	subtotalVenta := decimal.Zero
	consolidated := map[stockKey]decimal.Decimal{}

	for i, item := range items {
		result, deltas, motivo, err := s.priceLine(ctx, loader, item)
		if err != nil {
			return nil, decimal.Zero, nil, nil, err
		}
		if motivo != "" {
			loadErrors = append(loadErrors, domainventas.ItemErrorDetalle{Index: i, Motivo: motivo})
			continue
		}

		subtotalVenta = subtotalVenta.Add(result.subtotal)
		lines = append(lines, result)
		for _, d := range deltas {
			key := stockKey{TipoItem: d.TipoItem, ItemID: d.ItemID}
			consolidated[key] = consolidated[key].Add(d.Delta)
		}
	}

	if len(loadErrors) > 0 {
		return nil, decimal.Zero, nil, loadErrors, nil
	}

	movimientoInputs := make([]domainmovimientos.MovimientoInput, 0, len(consolidated))
	for key, cantidad := range consolidated {
		movimientoInputs = append(movimientoInputs, domainmovimientos.MovimientoInput{
			TipoItem:  domainstock.TipoItem(key.TipoItem),
			ItemID:    key.ItemID,
			Tipo:      domainmovimientos.TipoVenta,
			Ubicacion: domainstock.UbicacionVitrina,
			Cantidad:  cantidad,
		})
	}

	return lines, subtotalVenta, movimientoInputs, nil, nil
}

// priceEnvaseConFraganciaOrRecarga validates the fragancia and variante de
// envase (+ its modelo) a envase_con_fragancia/recarga line references, and
// computes its stock deltas: fragancia grams always come out, but the
// envase unit only comes out for envase_con_fragancia — a recarga doesn't
// consume envase stock, the customer brings their own.
func (s *Service) priceEnvaseConFraganciaOrRecarga(ctx context.Context, loader *entityLoader, item CreateVentaItemInput) (*envases.ModeloEnvase, []stockDelta, string, error) {
	if motivo, err := loader.checkFragancia(ctx, *item.FraganciaID); err != nil || motivo != "" {
		return nil, nil, motivo, err
	}
	modelo, motivo, err := loader.loadModeloByVariante(ctx, *item.VarianteEnvaseID)
	if err != nil || motivo != "" {
		return nil, nil, motivo, err
	}

	cantidadDec := decimal.NewFromInt32(item.Cantidad)
	gramosTotal := item.GramosFragancia.Mul(cantidadDec)
	deltas := []stockDelta{{TipoItem: string(domainstock.TipoItemFragancia), ItemID: *item.FraganciaID, Delta: gramosTotal.Neg()}}
	if item.TipoLinea == domainventas.TipoLineaEnvaseConFragancia {
		deltas = append(deltas, stockDelta{TipoItem: string(domainstock.TipoItemVarianteEnvase), ItemID: *item.VarianteEnvaseID, Delta: cantidadDec.Neg()})
	}
	return modelo, deltas, "", nil
}

// priceLine resolves and validates the catalog entities one line
// references, prices it, and returns the stock deltas it contributes. A
// non-empty motivo means the line failed entity validation (caller records
// it against this line's index); a non-nil error is a genuine
// infrastructure failure that aborts the whole request.
func (s *Service) priceLine(ctx context.Context, loader *entityLoader, item CreateVentaItemInput) (lineResult, []stockDelta, string, error) {
	pricingInput := domainventas.PricingInput{TipoLinea: item.TipoLinea, Cantidad: item.Cantidad}
	cantidadDec := decimal.NewFromInt32(item.Cantidad)
	var deltas []stockDelta

	switch item.TipoLinea {
	case domainventas.TipoLineaEnvaseConFragancia, domainventas.TipoLineaRecarga:
		modelo, lineDeltas, motivo, err := s.priceEnvaseConFraganciaOrRecarga(ctx, loader, item)
		if err != nil || motivo != "" {
			return lineResult{}, nil, motivo, err
		}
		pricingInput.ModeloEnvase = modelo
		deltas = lineDeltas

	case domainventas.TipoLineaEnvaseSolo:
		modelo, motivo, err := loader.loadModeloByVariante(ctx, *item.VarianteEnvaseID)
		if err != nil || motivo != "" {
			return lineResult{}, nil, motivo, err
		}
		pricingInput.ModeloEnvase = modelo
		deltas = append(deltas, stockDelta{TipoItem: string(domainstock.TipoItemVarianteEnvase), ItemID: *item.VarianteEnvaseID, Delta: cantidadDec.Neg()})

	case domainventas.TipoLineaProductoOtro:
		producto, motivo, err := loader.loadProducto(ctx, *item.ProductoID)
		if err != nil || motivo != "" {
			return lineResult{}, nil, motivo, err
		}
		pricingInput.Producto = producto
		deltas = append(deltas, stockDelta{TipoItem: string(domainstock.TipoItemProducto), ItemID: *item.ProductoID, Delta: cantidadDec.Neg()})
	}

	if item.FeromonaProductoID != nil {
		feromona, motivo, err := loader.loadFeromonaProducto(ctx, *item.FeromonaProductoID)
		if err != nil || motivo != "" {
			return lineResult{}, nil, motivo, err
		}
		pricingInput.FeromonaProducto = feromona
		deltas = append(deltas, stockDelta{TipoItem: string(domainstock.TipoItemProducto), ItemID: *item.FeromonaProductoID, Delta: cantidadDec.Neg()})
	}

	result, err := s.Pricing.Calculate(pricingInput)
	if err != nil {
		return lineResult{}, nil, "", internalErr(err)
	}

	return lineResult{item: item, precioUnitario: result.PrecioUnitario, subtotal: result.SubtotalLinea}, deltas, "", nil
}

// indexByStockKey rebuilds a stockKey->first-item-index map from the
// original request, used to attach an index to InventoryService's
// stock-insuficiente response (which has no positional information).
func indexByStockKey(items []CreateVentaItemInput) map[stockKey]int {
	result := map[stockKey]int{}
	record := func(tipoItem string, id *int64, index int) {
		if id == nil {
			return
		}
		key := stockKey{TipoItem: tipoItem, ItemID: *id}
		if _, ok := result[key]; !ok {
			result[key] = index
		}
	}
	for i, item := range items {
		record(string(domainstock.TipoItemFragancia), item.FraganciaID, i)
		record(string(domainstock.TipoItemVarianteEnvase), item.VarianteEnvaseID, i)
		record(string(domainstock.TipoItemProducto), item.ProductoID, i)
		record(string(domainstock.TipoItemProducto), item.FeromonaProductoID, i)
	}
	return result
}

func stockInsuficienteToDomainError(items []domainstock.StockInsuficienteItem, indexByKey map[stockKey]int) error {
	mapped := make([]domainventas.StockInsuficienteVentaItem, len(items))
	for i, it := range items {
		unidad := "unidades"
		if it.TipoItem == string(domainstock.TipoItemFragancia) {
			unidad = "gramos"
		}
		mapped[i] = domainventas.StockInsuficienteVentaItem{
			Index:      indexByKey[stockKey{TipoItem: it.TipoItem, ItemID: it.ItemID}],
			TipoItem:   it.TipoItem,
			ItemID:     it.ItemID,
			Nombre:     it.Nombre,
			Requerido:  it.Requerido,
			Disponible: it.Disponible,
			Unidad:     unidad,
		}
	}
	err := domainerrors.NewBusinessRule(
		"Stock insuficiente",
		"Uno o más ítems no tienen stock suficiente para completar la venta.",
	)
	err.Extra = domainventas.StockInsuficienteExtra{Items: mapped}
	return err
}

// entityLoader memoizes catalog lookups within one CreateVenta transaction
// so an entity referenced by several lines is fetched only once.
type entityLoader struct {
	fragancias       fraganciasrepo.Repository
	variantes        variantesenvaserepo.Repository
	modelos          modelosenvaserepo.Repository
	productos        productosrepo.Repository
	fraganciaCache   map[int64]struct{}
	modeloByVariante map[int64]*envases.ModeloEnvase
	productoCache    map[int64]*productos.Producto
}

func newEntityLoader(tx pgx.Tx) *entityLoader {
	return &entityLoader{
		fragancias:       fraganciasrepo.NewPostgres(tx),
		variantes:        variantesenvaserepo.NewPostgres(tx),
		modelos:          modelosenvaserepo.NewPostgres(tx),
		productos:        productosrepo.NewPostgres(tx),
		fraganciaCache:   map[int64]struct{}{},
		modeloByVariante: map[int64]*envases.ModeloEnvase{},
		productoCache:    map[int64]*productos.Producto{},
	}
}

// checkFragancia validates that id exists, isn't soft-deleted, and is
// active. Returns a non-empty motivo on validation failure, or a non-nil
// error only for genuine infrastructure failures.
func (l *entityLoader) checkFragancia(ctx context.Context, id int64) (string, error) {
	if _, ok := l.fraganciaCache[id]; ok {
		return "", nil
	}
	f, err := l.fragancias.GetByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, fraganciasrepo.ErrNotFound) {
			return "La fragancia indicada no existe.", nil
		}
		return "", internalErr(err)
	}
	if f.DeletedAt.Valid || !f.Activo {
		return "La fragancia indicada está eliminada o inactiva.", nil
	}
	l.fraganciaCache[id] = struct{}{}
	return "", nil
}

// loadModeloByVariante validates variante_envase_id, then resolves and
// validates its parent modelo_envase (variantes_envase doesn't carry the
// envase prices itself).
func (l *entityLoader) loadModeloByVariante(ctx context.Context, varianteID int64) (*envases.ModeloEnvase, string, error) {
	if modelo, ok := l.modeloByVariante[varianteID]; ok {
		return modelo, "", nil
	}

	variante, err := l.variantes.GetByIDIncludingDeleted(ctx, varianteID)
	if err != nil {
		if errors.Is(err, variantesenvaserepo.ErrNotFound) {
			return nil, "La variante de envase indicada no existe.", nil
		}
		return nil, "", internalErr(err)
	}
	if variante.DeletedAt.Valid || !variante.Activo {
		return nil, "La variante de envase indicada está eliminada o inactiva.", nil
	}

	modeloRow, err := l.modelos.GetByIDIncludingDeleted(ctx, variante.ModeloEnvaseID)
	if err != nil {
		if errors.Is(err, modelosenvaserepo.ErrNotFound) {
			return nil, "El modelo de envase de la variante indicada no existe.", nil
		}
		return nil, "", internalErr(err)
	}
	if modeloRow.DeletedAt.Valid || !modeloRow.Activo {
		return nil, "El modelo de envase de la variante indicada está eliminado o inactivo.", nil
	}

	modelo := &envases.ModeloEnvase{
		ID:                 modeloRow.ID,
		Tipo:               modeloRow.Tipo,
		TamanoOz:           modeloRow.TamanoOz,
		EquivGramos:        modeloRow.EquivGramos,
		PrecioSolo:         modeloRow.PrecioSolo,
		PrecioConFragancia: modeloRow.PrecioConFragancia,
		PrecioRecarga:      modeloRow.PrecioRecarga,
		Activo:             modeloRow.Activo,
	}
	l.modeloByVariante[varianteID] = modelo
	return modelo, "", nil
}

func (l *entityLoader) loadProducto(ctx context.Context, id int64) (*productos.Producto, string, error) {
	if p, ok := l.productoCache[id]; ok {
		return p, "", nil
	}
	row, err := l.productos.GetByIDIncludingDeleted(ctx, id)
	if err != nil {
		if errors.Is(err, productosrepo.ErrNotFound) {
			return nil, "El producto indicado no existe.", nil
		}
		return nil, "", internalErr(err)
	}
	if row.DeletedAt.Valid || !row.Activo {
		return nil, "El producto indicado está eliminado o inactivo.", nil
	}

	p := &productos.Producto{
		ID:        row.ID,
		SedeID:    row.SedeID,
		Nombre:    row.Nombre,
		Categoria: productos.Categoria(row.Categoria),
		Precio:    row.Precio,
		Activo:    row.Activo,
	}
	l.productoCache[id] = p
	return p, "", nil
}

// loadFeromonaProducto is loadProducto plus the extra categoria=feromona
// check the CHECK constraint doesn't enforce at the DB level.
func (l *entityLoader) loadFeromonaProducto(ctx context.Context, id int64) (*productos.Producto, string, error) {
	p, motivo, err := l.loadProducto(ctx, id)
	if err != nil || motivo != "" {
		return nil, motivo, err
	}
	if p.Categoria != productos.CategoriaFeromona {
		return nil, "feromona_producto_id no corresponde a un producto de categoría feromona.", nil
	}
	return p, "", nil
}

// loadVentaCompleta re-reads a venta with its items after commit, mapping
// generated rows to the domain Venta/VentaItem shape.
func (s *Service) loadVentaCompleta(ctx context.Context, id int64) (*domainventas.Venta, error) {
	row, err := s.Ventas.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ventasrepo.ErrNotFound) {
			return nil, notFoundErr()
		}
		return nil, internalErr(err)
	}
	itemRows, err := s.VentaItems.GetByVentaID(ctx, id)
	if err != nil {
		return nil, internalErr(err)
	}

	var observaciones *string
	if row.Observaciones.Valid {
		observaciones = &row.Observaciones.String
	}

	items := make([]domainventas.VentaItem, len(itemRows))
	for i, ir := range itemRows {
		items[i] = toDomainVentaItem(ir)
	}

	return &domainventas.Venta{
		ID:               row.ID,
		SedeID:           row.SedeID,
		UsuarioID:        row.UsuarioID,
		UsuarioNombre:    row.UsuarioNombre,
		MetodoPagoID:     row.MetodoPagoID,
		MetodoPagoNombre: row.MetodoPagoNombre,
		MetodoPagoCodigo: row.MetodoPagoCodigo,
		Subtotal:         row.Subtotal,
		DescuentoPct:     row.DescuentoPct,
		DescuentoMonto:   row.DescuentoMonto,
		Total:            row.Total,
		Observaciones:    observaciones,
		CreatedAt:        row.CreatedAt.Time,
		Items:            items,
	}, nil
}

func toDomainVentaItem(ir generated.GetVentaItemsByVentaIDRow) domainventas.VentaItem {
	var fraganciaID, varianteID, productoID, feromonaID *int64
	if ir.FraganciaID.Valid {
		v := ir.FraganciaID.Int64
		fraganciaID = &v
	}
	if ir.VarianteEnvaseID.Valid {
		v := ir.VarianteEnvaseID.Int64
		varianteID = &v
	}
	if ir.ProductoID.Valid {
		v := ir.ProductoID.Int64
		productoID = &v
	}
	if ir.FeromonaProductoID.Valid {
		v := ir.FeromonaProductoID.Int64
		feromonaID = &v
	}
	var gramos *decimal.Decimal
	if ir.GramosFragancia.Valid {
		g := ir.GramosFragancia.Decimal
		gramos = &g
	}

	var fraganciaNombre, productoNombre, feromonaNombre, varianteNombre *string
	if ir.FraganciaNombre.Valid {
		v := ir.FraganciaNombre.String
		fraganciaNombre = &v
	}
	if ir.ProductoNombre.Valid {
		v := ir.ProductoNombre.String
		productoNombre = &v
	}
	if ir.FeromonaNombre.Valid {
		v := ir.FeromonaNombre.String
		feromonaNombre = &v
	}
	if ir.ModeloEnvaseTipo.Valid && ir.VarianteEnvaseColor.Valid && ir.ModeloEnvaseTamanoOz.Valid {
		v := ir.ModeloEnvaseTipo.String + " " + ir.ModeloEnvaseTamanoOz.Decimal.String() + "oz " + ir.VarianteEnvaseColor.String
		varianteNombre = &v
	}

	return domainventas.VentaItem{
		ID:                   ir.ID,
		VentaID:              ir.VentaID,
		TipoLinea:            domainventas.TipoLinea(ir.TipoLinea),
		FraganciaID:          fraganciaID,
		FraganciaNombre:      fraganciaNombre,
		VarianteEnvaseID:     varianteID,
		VarianteEnvaseNombre: varianteNombre,
		ProductoID:           productoID,
		ProductoNombre:       productoNombre,
		FeromonaProductoID:   feromonaID,
		FeromonaNombre:       feromonaNombre,
		GramosFragancia:      gramos,
		Cantidad:             ir.Cantidad,
		PrecioUnitario:       ir.PrecioUnitario,
		Subtotal:             ir.Subtotal,
		CreatedAt:            ir.CreatedAt.Time,
	}
}
