package cuadres_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/consignaciones"
	cuadresrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/cuadres"
	pagoscajarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/pagos_caja"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	usecasecuadres "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
)

type testEnv struct {
	pool           *pgxpool.Pool
	service        *usecasecuadres.Service
	sedeID         int64
	adminID        int64
	vendedoraID    int64
	metodoEfectivo int64
	metodoNequi    int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE cuadres_caja, pagos_caja, consignaciones, ventas, venta_items,
		 metodos_pago, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	adminID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")
	vendedoraID := seedUsuario(t, pool, sedeID, "vendedora@test.local", "vendedora")
	metodoEfectivo := seedMetodoPago(t, pool, "Efectivo", "efectivo")
	metodoNequi := seedMetodoPago(t, pool, "Nequi", "nequi")

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.FixedZone("America/Bogota", -5*60*60)
	}

	service := usecasecuadres.NewService(
		pool,
		cuadresrepo.NewPostgres(pool),
		pagoscajarepo.NewPostgres(pool),
		consignaciones.NewPostgres(pool),
		usuarios.NewPostgres(pool),
		auditoria.NewPostgres(pool),
		loc,
	)

	return &testEnv{
		pool: pool, service: service, sedeID: sedeID,
		adminID: adminID, vendedoraID: vendedoraID,
		metodoEfectivo: metodoEfectivo, metodoNequi: metodoNequi,
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

// seedVenta inserts a minimal venta row directly (bypassing usecase/ventas,
// which this package doesn't depend on) so TotalesService has something to
// sum. subtotal == total, no discount.
func seedVenta(t *testing.T, env *testEnv, metodoPagoID int64, total string) {
	t.Helper()
	_, err := env.pool.Exec(context.Background(),
		`INSERT INTO ventas (sede_id, usuario_id, metodo_pago_id, subtotal, descuento_pct, descuento_monto, total)
		 VALUES ($1, $2, $3, $4::numeric, 0, 0, $4::numeric)`,
		env.sedeID, env.vendedoraID, metodoPagoID, total,
	)
	if err != nil {
		t.Fatalf("seeding venta: %v", err)
	}
}
