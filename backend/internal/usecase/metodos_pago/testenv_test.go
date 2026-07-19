package metodospago_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/metodos_pago"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/metodos_pago"
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

	// Truncate so no test sees rows left over by an earlier one sharing
	// this pool.
	if _, err := pool.Exec(context.Background(),
		`TRUNCATE ventas, metodos_pago, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	requesterID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")

	service := usecase.NewService(
		repo.NewPostgres(pool),
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

// seedVenta inserts a minimal venta row directly via SQL, used to exercise
// metodos_pago's conditional hard/soft delete rule.
func seedVenta(t *testing.T, pool *pgxpool.Pool, sedeID, usuarioID, metodoPagoID int64) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		`INSERT INTO ventas (sede_id, usuario_id, metodo_pago_id, subtotal, total)
		 VALUES ($1, $2, $3, 10000.00, 10000.00)`,
		sedeID, usuarioID, metodoPagoID,
	)
	if err != nil {
		t.Fatalf("seeding venta: %v", err)
	}
}
