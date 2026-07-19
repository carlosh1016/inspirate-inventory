package variantesenvase_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	modelosenvase "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/modelos_envase"
	stockactual "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/stock_actual"
	repo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/variantes_envase"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/variantes_envase"
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
		`TRUNCATE variantes_envase, modelos_envase, stock_actual, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	requesterID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")

	service := usecase.NewService(
		pool,
		repo.NewPostgres(pool),
		modelosenvase.NewPostgres(pool),
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

// seedModeloEnvase inserts a modelo_envase row directly via SQL so
// variantes_envase tests have a parent to reference without going through
// the modelos_envase usecase (a separate package).
func seedModeloEnvase(t *testing.T, pool *pgxpool.Pool, tipo, tamanoOz string, activo bool) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		`INSERT INTO modelos_envase (tipo, tamano_oz, equiv_gramos, precio_solo, precio_con_fragancia, precio_recarga, activo)
		 VALUES ($1, $2::numeric, 50.00, 10000.00, 25000.00, 15000.00, $3) RETURNING id`,
		tipo, tamanoOz, activo,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding modelo_envase: %v", err)
	}
	return id
}

// setStock overwrites the vitrina stock row cantidad for a variante_envase,
// used to exercise the delete-with-stock scenario without going through
// Tanda 3's (not-yet-built) movimientos usecase.
func setStock(t *testing.T, pool *pgxpool.Pool, itemID int64, cantidad string) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		`UPDATE stock_actual SET cantidad = $1::numeric
		 WHERE tipo_item = $2::tipo_item_enum AND item_id = $3 AND ubicacion = 'vitrina'::ubicacion_enum`,
		cantidad, stockactual.TipoItemVarianteEnvase, itemID,
	)
	if err != nil {
		t.Fatalf("setting stock: %v", err)
	}
}
