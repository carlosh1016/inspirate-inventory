package auditoria_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	auditoriarepo "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	usecaseauditoria "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auditoria"
)

type testEnv struct {
	pool       *pgxpool.Pool
	service    *usecaseauditoria.Service
	sedeID     int64
	usuarioID  int64
	usuario2ID int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	if _, err := pool.Exec(context.Background(),
		`TRUNCATE auditoria, usuarios, sedes RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	e := &testEnv{pool: pool}
	e.sedeID = seedSede(t, pool, "Sede Test")
	e.usuarioID = seedUsuario(t, pool, e.sedeID, "admin@test.local", "admin")
	e.usuario2ID = seedUsuario(t, pool, e.sedeID, "vendedora@test.local", "vendedora")
	e.service = usecaseauditoria.NewService(auditoriarepo.NewPostgres(pool))
	return e
}

func seedSede(t *testing.T, pool *pgxpool.Pool, nombre string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO sedes (nombre) VALUES ($1) RETURNING id`, nombre).Scan(&id); err != nil {
		t.Fatalf("seeding sede: %v", err)
	}
	return id
}

func seedUsuario(t *testing.T, pool *pgxpool.Pool, sedeID int64, correo, rol string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`INSERT INTO usuarios (sede_id, nombre_completo, correo, password_hash, rol)
		 VALUES ($1, 'Usuario Prueba', $2, 'x', $3) RETURNING id`,
		sedeID, correo, rol).Scan(&id); err != nil {
		t.Fatalf("seeding usuario: %v", err)
	}
	return id
}

// seedEvento inserts one audit row. usuarioID nil records an actor-less event;
// tabla/datos empty strings map to NULL.
func (e *testEnv) seedEvento(t *testing.T, usuarioID *int64, accion, tabla, datosDespues string, createdAt time.Time) {
	t.Helper()
	var tablaArg, datosArg any
	if tabla != "" {
		tablaArg = tabla
	}
	if datosDespues != "" {
		datosArg = datosDespues
	}
	if _, err := e.pool.Exec(context.Background(),
		`INSERT INTO auditoria (usuario_id, accion, tabla_afectada, datos_despues, ip, user_agent, created_at)
		 VALUES ($1, $2, $3, $4::jsonb, $5::inet, $6, $7::timestamptz)`,
		usuarioID, accion, tablaArg, datosArg, "192.168.1.10", "test-agent", createdAt,
	); err != nil {
		t.Fatalf("seeding evento: %v", err)
	}
}
