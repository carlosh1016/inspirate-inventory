package auth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/mailer"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/ratelimit"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	passwordresets "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/password_resets"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

// testPassword is the password every seeded usuario gets; individual tests
// exercise wrong/new passwords by passing literals straight to the usecase.
const testPassword = "password123"

type testEnv struct {
	pool    *pgxpool.Pool
	service *usecase.Service
	mailer  *mailer.MockMailer
	sedeID  int64
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	pool := requireTestPool(t)

	mockMailer := mailer.NewMock()
	service := usecase.NewService(
		usuarios.NewPostgres(pool),
		refreshtokens.NewPostgres(pool),
		passwordresets.NewPostgres(pool),
		auditoria.NewPostgres(pool),
		jwt.New("test-secret-for-integration-tests"),
		mockMailer,
		ratelimit.NewLoginLimiter(),
		ratelimit.NewPasswordResetLimiter(),
		24*time.Hour,
		10*time.Minute,
		720*time.Hour,
		8*time.Hour,
		"http://localhost:3000",
	)

	return &testEnv{
		pool:    pool,
		service: service,
		mailer:  mockMailer,
		sedeID:  seedSede(t, pool, "Sede "+t.Name()),
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

// sha256Hex mirrors the service's own (unexported) token hashing, so tests
// can seed refresh/reset tokens directly and then present the plaintext.
func sha256Hex(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
