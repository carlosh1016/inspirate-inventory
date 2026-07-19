package usuarios_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/usuarios"
)

// testPassword is the password every seeded usuario gets.
const testPassword = "password123"

type testEnv struct {
	pool    *pgxpool.Pool
	service *usecase.Service
	sedeID  int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	// CountActiveAdmins and List have no sede_id filter (matching the
	// original query spec), so every test needs a clean usuarios table —
	// otherwise admin-quorum and pagination assertions would be
	// contaminated by rows left over from earlier tests sharing this pool.
	if _, err := pool.Exec(context.Background(),
		`TRUNCATE usuarios, sedes, refresh_tokens, password_resets, auditoria RESTART IDENTITY CASCADE`,
	); err != nil {
		t.Fatalf("truncating tables between tests: %v", err)
	}

	service := usecase.NewService(
		usuarios.NewPostgres(pool),
		refreshtokens.NewPostgres(pool),
		auditoria.NewPostgres(pool),
	)

	return &testEnv{
		pool:    pool,
		service: service,
		sedeID:  seedSede(t, pool, "Sede "+t.Name()),
	}
}

// uniqueCorreo derives a correo from the test name: correo is globally
// unique (uq_usuarios_correo has no sede_id), and every test in this
// package shares one Postgres container/pool.
func uniqueCorreo(t *testing.T, suffix string) string {
	t.Helper()
	name := strings.ToLower(strings.NewReplacer("/", "-", " ", "-").Replace(t.Name()))
	return name + suffix + "@inspirate.co"
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

func seedUsuario(t *testing.T, pool *pgxpool.Pool, sedeID int64, correo, rol string, isActive bool) int64 {
	t.Helper()

	hash, err := domainauth.HashPassword(testPassword)
	if err != nil {
		t.Fatalf("hashing password: %v", err)
	}

	var id int64
	err = pool.QueryRow(context.Background(),
		`INSERT INTO usuarios (sede_id, nombre_completo, correo, password_hash, rol, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		sedeID, "Usuario de Prueba", correo, hash, rol, isActive,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seeding usuario: %v", err)
	}
	return id
}

// sha256Hex mirrors how the auth service hashes opaque tokens, so tests can
// seed refresh tokens directly.
func sha256Hex(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
