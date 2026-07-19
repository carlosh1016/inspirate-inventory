package stock_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	stockactualrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/stock"
)

type testEnv struct {
	pool    *pgxpool.Pool
	service *usecase.Service
	sedeID  int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE fragancias, variantes_envase, modelos_envase, productos, stock_actual, sedes RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")

	service := usecase.NewService(stockactualrepo.NewPostgres(pool))

	return &testEnv{pool: pool, service: service, sedeID: sedeID}
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

// seedFragancia inserts a fragancia (gramos_minimo fixed at 10.00, the only
// value the current tests need) and initializes its stock rows (both zero),
// returning its id.
func seedFragancia(t *testing.T, env *testEnv, nombre string, activo bool) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO fragancias (sede_id, nombre_comercial, genero, gramos_minimo, activo)
		 VALUES ($1, $2, 'masculina', 10.00, $3) RETURNING id`,
		env.sedeID, nombre, activo,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding fragancia: %v", err)
	}
	initStock(t, env, stockactualrepo.TipoItemFragancia, id)
	return id
}

// seedModeloEnvase inserts a modelo_envase, returning its id. No stock of
// its own — stock lives on its variantes.
func seedModeloEnvase(t *testing.T, env *testEnv, tipo, tamanoOz string) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO modelos_envase (tipo, tamano_oz, equiv_gramos, precio_solo, precio_con_fragancia, precio_recarga)
		 VALUES ($1, $2::numeric, 50.00, 10000.00, 25000.00, 15000.00) RETURNING id`,
		tipo, tamanoOz,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding modelo_envase: %v", err)
	}
	return id
}

// seedVarianteEnvase inserts a variante_envase and initializes its stock
// rows, returning its id.
func seedVarianteEnvase(t *testing.T, env *testEnv, modeloEnvaseID int64, color string, stockMinimo int32, activo bool) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO variantes_envase (modelo_envase_id, sede_id, color, stock_minimo, activo)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		modeloEnvaseID, env.sedeID, color, stockMinimo, activo,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding variante_envase: %v", err)
	}
	initStock(t, env, stockactualrepo.TipoItemVarianteEnvase, id)
	return id
}

// seedProducto inserts a producto and initializes its stock rows.
func seedProducto(t *testing.T, env *testEnv, nombre string, stockMinimo int32, activo bool) {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO productos (sede_id, nombre, categoria, precio, stock_minimo, activo)
		 VALUES ($1, $2, 'hogar', 20000.00, $3, $4) RETURNING id`,
		env.sedeID, nombre, stockMinimo, activo,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding producto: %v", err)
	}
	initStock(t, env, stockactualrepo.TipoItemProducto, id)
}

func initStock(t *testing.T, env *testEnv, tipoItem string, itemID int64) {
	t.Helper()
	if err := stockactualrepo.NewPostgres(env.pool).InitializeStock(context.Background(), env.sedeID, tipoItem, itemID); err != nil {
		t.Fatalf("initializing stock for %s %d: %v", tipoItem, itemID, err)
	}
}

// setVitrinaStock overwrites the vitrina stock row cantidad for an item.
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
