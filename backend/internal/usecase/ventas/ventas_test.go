package ventas_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	usecaseventas "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
)

func assertCode(t *testing.T, err error, code domainerrors.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var domainErr *domainerrors.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected *DomainError, got %T: %v", err, err)
	}
	if domainErr.Code != code {
		t.Fatalf("expected code %q, got %q (%v)", code, domainErr.Code, domainErr)
	}
}

func gramos(s string) *decimal.Decimal {
	d := decimal.RequireFromString(s)
	return &d
}

// --- 1. Venta simple con envase + fragancia ---

func TestCreateVentaEnvaseConFraganciaDescuentaStock(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	setVitrinaStock(t, env, "fragancia", fraganciaID, "50.00")
	setVitrinaStock(t, env, "variante_envase", varianteID, "10")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, GramosFragancia: gramos("11.00"), Cantidad: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Venta.Total.String() != "25000" {
		t.Fatalf("expected total=25000, got %s", out.Venta.Total)
	}

	var movCount int
	if err := env.pool.QueryRow(ctx, `SELECT COUNT(*) FROM movimientos_inventario WHERE venta_id = $1`, out.Venta.ID).Scan(&movCount); err != nil {
		t.Fatalf("counting movimientos: %v", err)
	}
	if movCount != 2 {
		t.Fatalf("expected 2 movimientos (fragancia + variante), got %d", movCount)
	}

	var fraganciaStock, varianteStock string
	if err := env.pool.QueryRow(ctx, `SELECT cantidad::text FROM stock_actual WHERE tipo_item = 'fragancia' AND item_id = $1 AND ubicacion = 'vitrina'`, fraganciaID).Scan(&fraganciaStock); err != nil {
		t.Fatalf("reading fragancia stock: %v", err)
	}
	if err := env.pool.QueryRow(ctx, `SELECT cantidad::text FROM stock_actual WHERE tipo_item = 'variante_envase' AND item_id = $1 AND ubicacion = 'vitrina'`, varianteID).Scan(&varianteStock); err != nil {
		t.Fatalf("reading variante stock: %v", err)
	}
	if fraganciaStock != "39.00" || varianteStock != "9.00" {
		t.Fatalf("expected fragancia=39.00 variante=9.00 after the sale, got fragancia=%s variante=%s", fraganciaStock, varianteStock)
	}
}

// --- 2. Venta con 4 tipos de línea mixtos ---

func TestCreateVentaCuatroTiposDeLinea(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	fraganciaA := seedFragancia(t, env, "Sauvage")
	modeloID := seedModeloEnvase(t, env)
	varianteA := seedVarianteEnvase(t, env, modeloID, "Rojo")
	varianteB := seedVarianteEnvase(t, env, modeloID, "Azul")
	productoID := seedProducto(t, env, "Vela de Vainilla")
	setVitrinaStock(t, env, "fragancia", fraganciaA, "100.00")
	setVitrinaStock(t, env, "variante_envase", varianteA, "10")
	setVitrinaStock(t, env, "variante_envase", varianteB, "10")
	setVitrinaStock(t, env, "producto", productoID, "10")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaA, VarianteEnvaseID: &varianteA, GramosFragancia: gramos("10.00"), Cantidad: 1},
			{TipoLinea: domainventas.TipoLineaRecarga, FraganciaID: &fraganciaA, VarianteEnvaseID: &varianteA, GramosFragancia: gramos("10.00"), Cantidad: 1},
			{TipoLinea: domainventas.TipoLineaEnvaseSolo, VarianteEnvaseID: &varianteB, Cantidad: 1},
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 25000 (envase_con_fragancia) + 15000 (recarga) + 10000 (envase_solo) + 20000 (producto_otro) = 70000
	if out.Venta.Subtotal.String() != "70000" {
		t.Fatalf("expected subtotal=70000, got %s", out.Venta.Subtotal)
	}
	if len(out.Venta.Items) != 4 {
		t.Fatalf("expected 4 venta_items, got %d", len(out.Venta.Items))
	}
}

// --- 3. Venta con feromona ---

func TestCreateVentaConFeromonaDescuentaSuStock(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	feromonaID := seedFeromona(t, env, "Feromona X")
	setVitrinaStock(t, env, "fragancia", fraganciaID, "50.00")
	setVitrinaStock(t, env, "variante_envase", varianteID, "10")
	setVitrinaStock(t, env, "producto", feromonaID, "5")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, FeromonaProductoID: &feromonaID, GramosFragancia: gramos("11.00"), Cantidad: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 25000 (envase con fragancia) + 1000 (feromona) = 26000
	if out.Venta.Subtotal.String() != "26000" {
		t.Fatalf("expected subtotal=26000, got %s", out.Venta.Subtotal)
	}

	var feromonaStock string
	if err := env.pool.QueryRow(ctx,
		`SELECT cantidad::text FROM stock_actual WHERE tipo_item = 'producto' AND item_id = $1 AND ubicacion = 'vitrina'`, feromonaID,
	).Scan(&feromonaStock); err != nil {
		t.Fatalf("reading feromona stock: %v", err)
	}
	if feromonaStock != "4.00" {
		t.Fatalf("expected feromona stock=4.00 (5-1), got %s", feromonaStock)
	}
}

// --- 4. Venta con recarga: no descuenta envase ---

func TestCreateVentaRecargaNoDescuentaEnvase(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	setVitrinaStock(t, env, "fragancia", fraganciaID, "50.00")
	setVitrinaStock(t, env, "variante_envase", varianteID, "10")

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaRecarga, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, GramosFragancia: gramos("11.00"), Cantidad: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var varianteStock string
	if err := env.pool.QueryRow(ctx,
		`SELECT cantidad::text FROM stock_actual WHERE tipo_item = 'variante_envase' AND item_id = $1 AND ubicacion = 'vitrina'`,
		varianteID,
	).Scan(&varianteStock); err != nil {
		t.Fatalf("reading variante stock: %v", err)
	}
	if varianteStock != "10.00" {
		t.Fatalf("expected variante_envase stock unchanged at 10.00, got %s", varianteStock)
	}
}

// --- 5. Dos líneas de la misma fragancia se consolidan ---

func TestCreateVentaConsolidaGramosDeLaMismaFragancia(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	setVitrinaStock(t, env, "fragancia", fraganciaID, "50.00")
	setVitrinaStock(t, env, "variante_envase", varianteID, "10")

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, GramosFragancia: gramos("11.00"), Cantidad: 1},
			{TipoLinea: domainventas.TipoLineaRecarga, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, GramosFragancia: gramos("11.00"), Cantidad: 1},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var fraganciaStock string
	if err := env.pool.QueryRow(ctx,
		`SELECT cantidad::text FROM stock_actual WHERE tipo_item = 'fragancia' AND item_id = $1 AND ubicacion = 'vitrina'`,
		fraganciaID,
	).Scan(&fraganciaStock); err != nil {
		t.Fatalf("reading fragancia stock: %v", err)
	}
	if fraganciaStock != "28.00" {
		t.Fatalf("expected fragancia stock=28.00 (50-22 consolidated), got %s", fraganciaStock)
	}

	var movCount int
	if err := env.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM movimientos_inventario WHERE tipo_item = 'fragancia' AND item_id = $1`, fraganciaID,
	).Scan(&movCount); err != nil {
		t.Fatalf("counting movimientos: %v", err)
	}
	if movCount != 1 {
		t.Fatalf("expected 1 consolidated movimiento for the fragancia, got %d", movCount)
	}
}

// --- 6/7/8. Descuento ---

func TestCreateVentaDescuento5Porciento(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Kit Regalo")
	setVitrinaStock(t, env, "producto", productoID, "10")

	// precio fijo del producto es 20000; para llegar a exactamente 50000
	// pedimos cantidad tal que 20000*cantidad = 50000 no es entero, así que
	// verificamos con dos líneas de producto (cantidad 2 + cantidad 0.5 no
	// aplica, cantidad es entero) — usamos 3 unidades de 20000 = 60000, que
	// también cae en el rango >=50000 y <100000 (5%).
	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 3},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Venta.Subtotal.String() != "60000" {
		t.Fatalf("expected subtotal=60000, got %s", out.Venta.Subtotal)
	}
	if out.Venta.DescuentoPct.String() != "5" || out.Venta.DescuentoMonto.String() != "3000" || out.Venta.Total.String() != "57000" {
		t.Fatalf("expected 5%% discount (3000) total=57000, got pct=%s monto=%s total=%s", out.Venta.DescuentoPct, out.Venta.DescuentoMonto, out.Venta.Total)
	}
}

func TestCreateVentaDescuento7Porciento(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Kit Regalo")
	setVitrinaStock(t, env, "producto", productoID, "10")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 5}, // 100000
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Venta.Subtotal.String() != "100000" {
		t.Fatalf("expected subtotal=100000, got %s", out.Venta.Subtotal)
	}
	if out.Venta.DescuentoPct.String() != "7" || out.Venta.DescuentoMonto.String() != "7000" || out.Venta.Total.String() != "93000" {
		t.Fatalf("expected 7%% discount (7000) total=93000, got pct=%s monto=%s total=%s", out.Venta.DescuentoPct, out.Venta.DescuentoMonto, out.Venta.Total)
	}
}

func TestCreateVentaSinDescuento(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 2}, // 40000 < 49999 threshold isn't hit either
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Venta.DescuentoMonto.IsZero() || out.Venta.Total.String() != out.Venta.Subtotal.String() {
		t.Fatalf("expected no discount below 50000, got %+v", out.Venta)
	}
}

// --- 9/10. Stock insuficiente ---

func TestCreateVentaStockInsuficienteFragancia(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	setVitrinaStock(t, env, "fragancia", fraganciaID, "5.00")
	setVitrinaStock(t, env, "variante_envase", varianteID, "10")

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, GramosFragancia: gramos("20.00"), Cantidad: 1},
		},
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)

	var domainErr *domainerrors.DomainError
	errors.As(err, &domainErr)
	extra, ok := domainErr.Extra.(domainventas.StockInsuficienteExtra)
	if !ok || len(extra.Items) != 1 {
		t.Fatalf("expected StockInsuficienteExtra with 1 item, got %#v", domainErr.Extra)
	}
	if extra.Items[0].Index != 0 || extra.Items[0].Unidad != "gramos" {
		t.Fatalf("expected index=0 unidad=gramos, got %+v", extra.Items[0])
	}
}

func TestCreateVentaStockInsuficienteDosItems(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteA := seedVarianteEnvase(t, env, modeloID, "Morado")
	varianteB := seedVarianteEnvase(t, env, modeloID, "Verde")
	setVitrinaStock(t, env, "fragancia", fraganciaID, "5.00")
	setVitrinaStock(t, env, "variante_envase", varianteA, "10") // sufficient — only the fragancia portion of item 0 should fail
	setVitrinaStock(t, env, "variante_envase", varianteB, "0")

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteA, GramosFragancia: gramos("20.00"), Cantidad: 1},
			{TipoLinea: domainventas.TipoLineaEnvaseSolo, VarianteEnvaseID: &varianteB, Cantidad: 1},
		},
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	var domainErr *domainerrors.DomainError
	errors.As(err, &domainErr)
	extra, ok := domainErr.Extra.(domainventas.StockInsuficienteExtra)
	if !ok || len(extra.Items) != 2 {
		t.Fatalf("expected StockInsuficienteExtra with 2 items, got %#v", domainErr.Extra)
	}
}

// --- 11/12/13. Ítems inválidos ---

func TestCreateVentaFraganciaSoftDeletedFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	if _, err := env.pool.Exec(ctx, `UPDATE fragancias SET deleted_at = NOW() WHERE id = $1`, fraganciaID); err != nil {
		t.Fatalf("soft-deleting fragancia: %v", err)
	}

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, GramosFragancia: gramos("11.00"), Cantidad: 1},
		},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCreateVentaVarianteInactivaFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	if _, err := env.pool.Exec(ctx, `UPDATE variantes_envase SET activo = false WHERE id = $1`, varianteID); err != nil {
		t.Fatalf("deactivating variante: %v", err)
	}

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseSolo, VarianteEnvaseID: &varianteID, Cantidad: 1},
		},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

func TestCreateVentaFeromonaCategoriaInvalidaFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	fraganciaID := seedFragancia(t, env, "Bleu de Chanel")
	modeloID := seedModeloEnvase(t, env)
	varianteID := seedVarianteEnvase(t, env, modeloID, "Morado")
	noFeromona := seedProducto(t, env, "Vela") // categoria "hogar", not "feromona"
	setVitrinaStock(t, env, "fragancia", fraganciaID, "50.00")
	setVitrinaStock(t, env, "variante_envase", varianteID, "10")

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, FraganciaID: &fraganciaID, VarianteEnvaseID: &varianteID, FeromonaProductoID: &noFeromona, GramosFragancia: gramos("11.00"), Cantidad: 1},
		},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

// --- 14/15. Idempotency ---

func TestCreateVentaIdempotencyKeyRepetidoDevuelveMismaRespuesta(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	key := "idem-key-1"
	in := usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1},
		},
		IdempotencyKey:  &key,
		RequestBodyHash: "hash-abc",
	}

	first, err := env.service.CreateVenta(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
	if first.FromIdempotency {
		t.Fatal("expected the first call to not be from idempotency cache")
	}

	// The real response caching is owned by the HTTP handler (it's the only
	// layer that can produce the real client-facing response), so here we
	// simulate what it does after a successful CreateVenta call.
	cachedBody := []byte(`{"data":{"id":` + strconv.FormatInt(first.Venta.ID, 10) + `}}`)
	if err := env.service.StoreIdempotencyResponse(ctx, key, env.vendedoraID, in.RequestBodyHash, http.StatusCreated, cachedBody); err != nil {
		t.Fatalf("storing idempotency response: %v", err)
	}

	second, err := env.service.CreateVenta(ctx, in)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if !second.FromIdempotency {
		t.Fatal("expected the second call to be served from the idempotency cache")
	}
	// response_body is stored as JSONB, so Postgres re-serializes it (e.g.
	// adds a space after ":") — compare parsed JSON, not raw bytes.
	var wantJSON, gotJSON map[string]any
	if err := json.Unmarshal(cachedBody, &wantJSON); err != nil {
		t.Fatalf("unmarshaling expected body: %v", err)
	}
	if err := json.Unmarshal(second.CachedBody, &gotJSON); err != nil {
		t.Fatalf("unmarshaling cached body: %v", err)
	}
	if second.CachedStatus != http.StatusCreated || !reflect.DeepEqual(wantJSON, gotJSON) {
		t.Fatalf("expected the cached response to be replayed, got status=%d body=%s", second.CachedStatus, second.CachedBody)
	}

	var ventaCount int
	if err := env.pool.QueryRow(ctx, `SELECT COUNT(*) FROM ventas`).Scan(&ventaCount); err != nil {
		t.Fatalf("counting ventas: %v", err)
	}
	if ventaCount != 1 {
		t.Fatalf("expected exactly 1 venta after repeating the idempotency key, got %d", ventaCount)
	}

	var stock string
	if err := env.pool.QueryRow(ctx,
		`SELECT cantidad::text FROM stock_actual WHERE tipo_item = 'producto' AND item_id = $1 AND ubicacion = 'vitrina'`, productoID,
	).Scan(&stock); err != nil {
		t.Fatalf("reading producto stock: %v", err)
	}
	if stock != "9.00" {
		t.Fatalf("expected stock deducted only once (9.00), got %s", stock)
	}
}

func TestCreateVentaIdempotencyKeyConDistintoBodyFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	key := "idem-key-2"
	first := usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1},
		},
		IdempotencyKey:  &key,
		RequestBodyHash: "hash-original",
	}
	firstOut, err := env.service.CreateVenta(ctx, first)
	if err != nil {
		t.Fatalf("unexpected error on first call: %v", err)
	}
	if err := env.service.StoreIdempotencyResponse(ctx, key, env.vendedoraID, first.RequestBodyHash, http.StatusCreated, []byte(`{"data":{"id":`+strconv.FormatInt(firstOut.Venta.ID, 10)+`}}`)); err != nil {
		t.Fatalf("storing idempotency response: %v", err)
	}

	second := first
	second.RequestBodyHash = "hash-diferente"
	_, err = env.service.CreateVenta(ctx, second)
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	assertCode(t, err, domainerrors.CodeConflict)
}

// --- 16/17. Ítems mal formados ---

func TestCreateVentaItemMalFormadoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaEnvaseConFragancia, Cantidad: 1}, // missing fragancia_id/variante_envase_id/gramos_fragancia
		},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)

	var domainErr *domainerrors.DomainError
	errors.As(err, &domainErr)
	extra, ok := domainErr.Extra.(domainventas.ItemErrorExtra)
	if !ok || len(extra.Items) != 1 || extra.Items[0].Index != 0 {
		t.Fatalf("expected ItemErrorExtra with 1 item at index 0, got %#v", domainErr.Extra)
	}
}

func TestCreateVentaSinItemsFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: nil,
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

// --- 18. Método de pago inactivo ---

func TestCreateVentaMetodoPagoInactivoFalla(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")
	if _, err := env.pool.Exec(ctx, `UPDATE metodos_pago SET activo = false WHERE id = $1`, env.metodoPagoID); err != nil {
		t.Fatalf("deactivating metodo_pago: %v", err)
	}

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{
			{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1},
		},
	})
	if err == nil {
		t.Fatal("expected a validation error, got nil")
	}
	assertCode(t, err, domainerrors.CodeValidation)
}

// --- List / Get / Update ---

func TestListAdminVeTodasVendedoraSoloLasSuyas(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	if _, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	}); err != nil {
		t.Fatalf("seeding venta de vendedora: %v", err)
	}
	if _, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.adminID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	}); err != nil {
		t.Fatalf("seeding venta de admin: %v", err)
	}

	all, err := env.service.List(ctx, usecaseventas.ListInput{SedeID: env.sedeID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if all.Total != 2 {
		t.Fatalf("expected admin view to see 2 ventas, got %d", all.Total)
	}

	onlyVendedora, err := env.service.List(ctx, usecaseventas.ListInput{SedeID: env.sedeID, UsuarioID: env.vendedoraID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if onlyVendedora.Total != 1 {
		t.Fatalf("expected vendedora-scoped view to see 1 venta, got %d", onlyVendedora.Total)
	}
}

func TestGetVentaDesconocidaNotFound(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Get(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected not found error, got nil")
	}
	assertCode(t, err, domainerrors.CodeNotFound)
}

func TestUpdateObservacionesAuditado(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	})
	if err != nil {
		t.Fatalf("seeding venta: %v", err)
	}

	obs := "Cliente pidió factura"
	updated, err := env.service.Update(ctx, usecaseventas.UpdateInput{
		TargetID: out.Venta.ID, Observaciones: &obs, RequesterID: env.adminID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Observaciones == nil || *updated.Observaciones != obs {
		t.Fatalf("expected observaciones=%q, got %+v", obs, updated.Observaciones)
	}

	var auditCount int
	if err := env.pool.QueryRow(ctx, `SELECT COUNT(*) FROM auditoria WHERE accion = 'venta_observaciones_editadas'`).Scan(&auditCount); err != nil {
		t.Fatalf("counting auditoria: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("expected 1 auditoria entry, got %d", auditCount)
	}
}

// --- ResumenHoy ---

func TestResumenHoyConVentas(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	if _, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	}); err != nil {
		t.Fatalf("seeding venta: %v", err)
	}

	resumen, err := env.service.ResumenHoy(ctx, env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resumen.Resumen.VentasCount != 1 {
		t.Fatalf("expected ventas_count=1, got %d", resumen.Resumen.VentasCount)
	}
	if resumen.Resumen.TotalDia.String() != "20000" {
		t.Fatalf("expected total_dia=20000, got %s", resumen.Resumen.TotalDia)
	}
}

func TestResumenHoySinVentas(t *testing.T) {
	env := newTestEnv(t)

	resumen, err := env.service.ResumenHoy(context.Background(), env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resumen.Resumen.VentasCount != 0 || !resumen.Resumen.TotalDia.IsZero() {
		t.Fatalf("expected zero ventas today, got %+v", resumen.Resumen)
	}
}

func TestResumenHoyExcluyeVentasDeAyer(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	})
	if err != nil {
		t.Fatalf("seeding venta: %v", err)
	}
	if _, err := env.pool.Exec(ctx, `UPDATE ventas SET created_at = NOW() - INTERVAL '1 day' WHERE id = $1`, out.Venta.ID); err != nil {
		t.Fatalf("backdating venta: %v", err)
	}

	resumen, err := env.service.ResumenHoy(ctx, env.sedeID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resumen.Resumen.VentasCount != 0 {
		t.Fatalf("expected yesterday's venta excluded from today's resumen, got ventas_count=%d", resumen.Resumen.VentasCount)
	}
}

// --- Integración con Tanda 5: CajaStatusService ---

func TestCreateVentaBloqueaSiCuadreEstaCerrado(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	env.cajaStatus.Err = domainerrors.NewBusinessRule(
		"Cuadre de caja cerrado",
		"El cuadre de caja del día está cerrado. No se pueden registrar más ventas para hoy.",
	)

	_, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	})
	if err == nil {
		t.Fatal("expected a business_rule error, got nil")
	}
	assertCode(t, err, domainerrors.CodeBusinessRule)

	var count int
	if err := env.pool.QueryRow(ctx, `SELECT COUNT(*) FROM ventas`).Scan(&count); err != nil {
		t.Fatalf("counting ventas: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no venta to have been created, got %d", count)
	}
}

func TestCreateVentaPermiteSiCuadreAbiertoOSinCuadre(t *testing.T) {
	env := newTestEnv(t)
	ctx := context.Background()
	productoID := seedProducto(t, env, "Vela")
	setVitrinaStock(t, env, "producto", productoID, "10")

	// env.cajaStatus.Err is nil by default (no cuadre / cuadre abierto).
	out, err := env.service.CreateVenta(ctx, usecaseventas.CreateVentaInput{
		SedeID: env.sedeID, UsuarioID: env.vendedoraID, MetodoPagoID: env.metodoPagoID,
		Items: []usecaseventas.CreateVentaItemInput{{TipoLinea: domainventas.TipoLineaProductoOtro, ProductoID: &productoID, Cantidad: 1}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Venta == nil {
		t.Fatal("expected the venta to be created")
	}
}
