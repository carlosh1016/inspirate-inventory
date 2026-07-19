package modelosenvase_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
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
		`TRUNCATE modelos_envase, variantes_envase, stock_actual, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
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

// seedVarianteEnvase inserts a variante_envase row directly via SQL (bypassing
// the variantes_envase usecase, which lives in a separate package) so
// modelos_envase's delete-blocking rule has something to collide with.
func seedVarianteEnvase(t *testing.T, pool *pgxpool.Pool, modeloEnvaseID, sedeID int64, color string) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		`INSERT INTO variantes_envase (modelo_envase_id, sede_id, color) VALUES ($1, $2, $3) RETURNING id`,
		modeloEnvaseID, sedeID, color,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding variante_envase: %v", err)
	}
	return id
}
