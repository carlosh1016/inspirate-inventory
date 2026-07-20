package sesiones_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	sesionesrepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/sesiones"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	usecasesesiones "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/sesiones"
)

type testEnv struct {
	pool        *pgxpool.Pool
	service     *usecasesesiones.Service
	sedeID      int64
	vendedoraID int64
	vendedora2  int64
	adminID     int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE sesiones_laborales, usuarios, sedes, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	sedeID := seedSede(t, pool, "Sede Test")
	adminID := seedUsuario(t, pool, sedeID, "admin@test.local", "admin")
	vendedoraID := seedUsuario(t, pool, sedeID, "vendedora@test.local", "vendedora")
	vendedora2 := seedUsuario(t, pool, sedeID, "vendedora2@test.local", "vendedora")

	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		loc = time.FixedZone("America/Bogota", -5*60*60)
	}

	service := usecasesesiones.NewService(
		sesionesrepo.NewPostgres(pool),
		usuarios.NewPostgres(pool),
		auditoria.NewPostgres(pool),
		loc,
	)

	return &testEnv{
		pool: pool, service: service, sedeID: sedeID,
		adminID: adminID, vendedoraID: vendedoraID, vendedora2: vendedora2,
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
