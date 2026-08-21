package movimientos_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	fraganciasrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	movimientosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/movimientos"
	productosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/productos"
	stockactualrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	variantesenvaserepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

type testEnv struct {
	pool        *pgxpool.Pool
	service     *usecase.Service
	sedeID      int64
	requesterID int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE movimientos_inventario, fragancias, variantes_envase, modelos_envase, productos,
		 stock_actual, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	requesterID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")

	service := usecase.NewService(
		pool,
		movimientosrepo.NewPostgres(pool),
		stockactualrepo.NewPostgres(pool),
		fraganciasrepo.NewPostgres(pool),
		variantesenvaserepo.NewPostgres(pool),
		productosrepo.NewPostgres(pool),
		auditoria.NewPostgres(pool),
	)

	return &testEnv{pool: pool, service: service, sedeID: sedeID, requesterID: requesterID}
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

// seedFragancia inserts a fragancia (nombre_comercial fixed at "Bleu de
// Chanel", gramos_minimo at 10.00 — the only values the current tests need)
// and initializes its stock rows (both zero), returning its id.
func seedFragancia(t *testing.T, env *testEnv) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO fragancias (sede_id, nombre_comercial, genero, gramos_minimo, numero_genero)
		 VALUES ($1, 'Bleu de Chanel', 'masculina', 10.00,
		   (SELECT COALESCE(MAX(numero_genero), 0) + 1 FROM fragancias WHERE sede_id = $1 AND genero = 'masculina'))
		 RETURNING id`,
		env.sedeID,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding fragancia: %v", err)
	}
	initStock(t, env, stockactualrepo.TipoItemFragancia, id)
	return id
}

// seedProducto inserts a producto and initializes its stock rows, returning
// its id.
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
	initStock(t, env, stockactualrepo.TipoItemProducto, id)
	return id
}

func initStock(t *testing.T, env *testEnv, tipoItem string, itemID int64) {
	t.Helper()
	if err := stockactualrepo.NewPostgres(env.pool).InitializeStock(context.Background(), env.sedeID, tipoItem, itemID); err != nil {
		t.Fatalf("initializing stock for %s %d: %v", tipoItem, itemID, err)
	}
}

// setStock overwrites the stock_actual cantidad for one fragancia/ubicacion,
// used to exercise scenarios without going through EntradaMercancia first.
func setStock(t *testing.T, env *testEnv, itemID int64, ubicacion, cantidad string) {
	t.Helper()
	_, err := env.pool.Exec(context.Background(),
		`UPDATE stock_actual SET cantidad = $1::numeric
		 WHERE tipo_item = 'fragancia'::tipo_item_enum AND item_id = $2 AND ubicacion = $3::ubicacion_enum`,
		cantidad, itemID, ubicacion,
	)
	if err != nil {
		t.Fatalf("setting stock: %v", err)
	}
}
