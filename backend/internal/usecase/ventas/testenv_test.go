package ventas_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	fraganciasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	idempotencykeysrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/idempotency_keys"
	metodospagorepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
	modelosenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	productosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactualrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	variantesenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	ventaitemsrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/venta_items"
	ventasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/ventas"
	usecasecuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
	usecasemovimientos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
	usecaseventas "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
)

type testEnv struct {
	pool         *pgxpool.Pool
	service      *usecaseventas.Service
	cajaStatus   *mockCajaStatusService
	sedeID       int64
	adminID      int64
	vendedoraID  int64
	metodoPagoID int64
}

// mockCajaStatusService is a stand-in for usecase/cuadres.CajaStatusService
// so ventas tests don't need a real cuadres_caja row: Err is nil by
// default (venta always allowed), and tests that need to exercise the
// block set it to a *domainerrors.DomainError before calling CreateVenta.
type mockCajaStatusService struct {
	Err error
}

var _ usecasecuadres.CajaStatusService = (*mockCajaStatusService)(nil)

func (m *mockCajaStatusService) VerificarPuedeRegistrarVenta(context.Context, int64, time.Time) error {
	return m.Err
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE ventas, venta_items, idempotency_keys, movimientos_inventario,
		 fragancias, variantes_envase, modelos_envase, productos, metodos_pago,
		 stock_actual, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	adminID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")
	vendedoraID := seedUsuario(t, pool, sedeID, "vendedora@test.local", "vendedora")
	metodoPagoID := seedMetodoPago(t, pool, "Efectivo", "efectivo")

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.FixedZone("America/Bogota", -5*60*60)
	}

	stockActualRepo := stockactualrepo.NewPostgres(pool)
	fraganciasRepo := fraganciasrepo.NewPostgres(pool)
	variantesEnvaseRepo := variantesenvaserepo.NewPostgres(pool)
	modelosEnvaseRepo := modelosenvaserepo.NewPostgres(pool)
	productosRepo := productosrepo.NewPostgres(pool)
	metodosPagoRepo := metodospagorepo.NewPostgres(pool)
	movimientosRepo := movimientosrepo.NewPostgres(pool)
	auditoriaRepo := auditoria.NewPostgres(pool)

	movimientosService := usecasemovimientos.NewService(
		pool, movimientosRepo, stockActualRepo, fraganciasRepo, variantesEnvaseRepo, productosRepo, auditoriaRepo,
	)

	cajaStatus := &mockCajaStatusService{}

	service := usecaseventas.NewService(
		pool,
		ventasrepo.NewPostgres(pool),
		ventaitemsrepo.NewPostgres(pool),
		idempotencykeysrepo.NewPostgres(pool),
		movimientosService,
		movimientosRepo,
		fraganciasRepo,
		variantesEnvaseRepo,
		modelosEnvaseRepo,
		productosRepo,
		metodosPagoRepo,
		auditoriaRepo,
		usecaseventas.NewPricingService(),
		usecaseventas.NewDiscountService(),
		cajaStatus,
		loc,
	)

	return &testEnv{
		pool: pool, service: service, cajaStatus: cajaStatus, sedeID: sedeID,
		adminID: adminID, vendedoraID: vendedoraID, metodoPagoID: metodoPagoID,
	}
}

func seedSede(t *testing.T, pool *pgxpool.Pool, nombre string) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(), `INSERT INTO sedes (nombre) VALUES ($1) RETURNING id`, nombre).Scan(&id)
	if err != nil {
		t.Fatalf("seeding sede: %v", err)
	}
	return id
}

func seedUsuario(t *testing.T, pool *pgxpool.Pool, sedeID int64, correo, rol string) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		`INSERT INTO usuarios (sede_id, nombre_completo, correo, password_hash, rol) VALUES ($1, 'Usuario Prueba', $2, 'x', $3) RETURNING id`,
		sedeID, correo, rol,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding usuario: %v", err)
	}
	return id
}

func seedMetodoPago(t *testing.T, pool *pgxpool.Pool, nombre, codigo string) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		`INSERT INTO metodos_pago (nombre, codigo) VALUES ($1, $2) RETURNING id`, nombre, codigo,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding metodo_pago: %v", err)
	}
	return id
}

// seedFragancia inserts a fragancia (gramos_minimo fixed at 10.00) and
// initializes its stock rows (both zero), returning its id.
func seedFragancia(t *testing.T, env *testEnv, nombre string) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO fragancias (sede_id, nombre_comercial, genero, gramos_minimo)
		 VALUES ($1, $2, 'masculina', 10.00) RETURNING id`,
		env.sedeID, nombre,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding fragancia: %v", err)
	}
	initStock(t, env, "fragancia", id)
	return id
}

// seedModeloEnvase inserts a modelo_envase (tipo "Spray", tamano_oz 1.00 —
// the only combination the current tests need) with fixed prices
// (solo=10000, con_fragancia=25000, recarga=15000), returning its id.
func seedModeloEnvase(t *testing.T, env *testEnv) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO modelos_envase (tipo, tamano_oz, equiv_gramos, precio_solo, precio_con_fragancia, precio_recarga)
		 VALUES ('Spray', 1.00, 50.00, 10000.00, 25000.00, 15000.00) RETURNING id`,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding modelo_envase: %v", err)
	}
	return id
}

// seedVarianteEnvase inserts a variante_envase and initializes its stock
// rows, returning its id.
func seedVarianteEnvase(t *testing.T, env *testEnv, modeloEnvaseID int64, color string) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO variantes_envase (modelo_envase_id, sede_id, color) VALUES ($1, $2, $3) RETURNING id`,
		modeloEnvaseID, env.sedeID, color,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding variante_envase: %v", err)
	}
	initStock(t, env, "variante_envase", id)
	return id
}

// seedProducto inserts a producto (categoria "hogar" by default, precio
// 20000) and initializes its stock rows, returning its id.
func seedProducto(t *testing.T, env *testEnv, nombre string) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO productos (sede_id, nombre, categoria, precio, stock_minimo)
		 VALUES ($1, $2, 'hogar', 20000.00, 0) RETURNING id`,
		env.sedeID, nombre,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding producto: %v", err)
	}
	initStock(t, env, "producto", id)
	return id
}

// seedFeromona inserts a producto with categoria='feromona' (precio 1000)
// and initializes its stock rows, returning its id.
func seedFeromona(t *testing.T, env *testEnv, nombre string) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO productos (sede_id, nombre, categoria, precio, stock_minimo)
		 VALUES ($1, $2, 'feromona', 1000.00, 0) RETURNING id`,
		env.sedeID, nombre,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding feromona: %v", err)
	}
	initStock(t, env, "producto", id)
	return id
}

func initStock(t *testing.T, env *testEnv, tipoItem string, itemID int64) {
	t.Helper()
	if err := stockactualrepo.NewPostgres(env.pool).InitializeStock(context.Background(), env.sedeID, tipoItem, itemID); err != nil {
		t.Fatalf("initializing stock for %s %d: %v", tipoItem, itemID, err)
	}
}

// setVitrinaStock overwrites the vitrina stock row cantidad for one item.
func setVitrinaStock(t *testing.T, env *testEnv, tipoItem string, itemID int64, cantidad string) {
	t.Helper()
	_, err := env.pool.Exec(context.Background(),
		`UPDATE stock_actual SET cantidad = $1::numeric
		 WHERE tipo_item = $2::tipo_item_enum AND item_id = $3 AND ubicacion = 'vitrina'::ubicacion_enum`,
		cantidad, tipoItem, itemID,
	)
	if err != nil {
		t.Fatalf("setting vitrina stock: %v", err)
	}
}
