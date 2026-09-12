package conteos_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	conteosrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/conteos"
	usecaseconteos "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/conteos"
)

type testEnv struct {
	pool    *pgxpool.Pool
	service *usecaseconteos.Service
	loc     *time.Location
	sedeID  int64
	adminID int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE conteos_inventario, conteo_inventario_items, fragancias, stock_actual, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	adminID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.FixedZone("America/Bogota", -5*60*60)
	}

	service := usecaseconteos.NewService(
		pool,
		conteosrepo.NewPostgres(pool),
		auditoria.NewPostgres(pool),
		loc,
	)

	return &testEnv{pool: pool, service: service, loc: loc, sedeID: sedeID, adminID: adminID}
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

// seedFragancia inserts an active fragancia for env.sedeID plus its
// stock_actual row (all in bodega, vitrina left at zero) with the given
// total grams.
func seedFragancia(t *testing.T, env *testEnv, nombre, stockGramos string) int64 {
	t.Helper()
	var id int64
	err := env.pool.QueryRow(context.Background(),
		`INSERT INTO fragancias (sede_id, nombre_comercial, genero, numero_genero, gramos_minimo)
		 VALUES ($1, $2, 'femenina', nextval('fragancias_id_seq'), 0) RETURNING id`,
		env.sedeID, nombre,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding fragancia: %v", err)
	}
	_, err = env.pool.Exec(context.Background(),
		`INSERT INTO stock_actual (sede_id, tipo_item, item_id, ubicacion, cantidad) VALUES
		 ($1, 'fragancia', $2, 'bodega', $3::numeric),
		 ($1, 'fragancia', $2, 'vitrina', 0)`,
		env.sedeID, id, stockGramos,
	)
	if err != nil {
		t.Fatalf("seeding stock_actual: %v", err)
	}
	return id
}

// mesAnterior/mesActual/mesSiguiente return the first day of the month
// relative to now, in env's timezone — matching how Service resolves
// "el mes actual" when Periodo is nil.
func mesActual(env *testEnv) time.Time {
	now := time.Now().In(env.loc)
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, env.loc)
}

func mesAnterior(env *testEnv) time.Time {
	return mesActual(env).AddDate(0, -1, 0)
}

// seedMovimiento inserts a movimientos_inventario row directly (bypassing
// usecase/movimientos, which this package doesn't depend on) so
// GenerarItemsConteo's ledger-based saldo_inicial reconstruction has
// something to read. cantidad carries its sign like the real engine does
// (positive for entradas, negative for salidas).
func seedMovimiento(t *testing.T, env *testEnv, fraganciaID int64, tipo, cantidad string) {
	t.Helper()
	_, err := env.pool.Exec(context.Background(),
		`INSERT INTO movimientos_inventario (sede_id, usuario_id, tipo_item, item_id, tipo, ubicacion, cantidad, stock_anterior, stock_posterior, motivo)
		 VALUES ($1, $2, 'fragancia', $3, $4::tipo_movimiento_enum, 'bodega', $5::numeric, 0, 0, 'seed de prueba')`,
		env.sedeID, env.adminID, fraganciaID, tipo, cantidad,
	)
	if err != nil {
		t.Fatalf("seeding movimiento: %v", err)
	}
}
