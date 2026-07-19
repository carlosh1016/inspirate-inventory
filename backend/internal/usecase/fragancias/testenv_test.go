package fragancias_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/fragancias"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
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
	// this pool (fragancias' uniqueness/list queries are sede-scoped, but
	// the sede/usuario ids themselves would otherwise keep growing).
	if _, err := pool.Exec(context.Background(),
		`TRUNCATE fragancias, stock_actual, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	requesterID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")

	service := usecase.NewService(
		pool,
		repo.NewPostgres(pool),
		stockactual.NewPostgres(pool),
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

// setStock overwrites the vitrina stock row cantidad for a fragancia, used
// to exercise delete-with-stock and stock_bajo scenarios without going
// through Tanda 3's (not-yet-built) movimientos usecase.
func setStock(t *testing.T, pool *pgxpool.Pool, itemID int64, cantidad string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		`UPDATE stock_actual SET cantidad = $1::numeric
		 WHERE tipo_item = $2::tipo_item_enum AND item_id = $3 AND ubicacion = 'vitrina'::ubicacion_enum`,
		cantidad, stockactual.TipoItemFragancia, itemID,
	)
	if err != nil {
		t.Fatalf("setting stock: %v", err)
	}
}
