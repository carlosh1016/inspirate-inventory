package auth_test

import (
	"context"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

func TestLoginRateLimited(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	in := usecase.LoginInput{Correo: correo, Password: testPassword, IP: "10.0.0.1", UserAgent: "ua"}
	for i := 0; i < 5; i++ {
		if _, err := env.service.Login(context.Background(), in); err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", i+1, err)
		}
	}

	_, err := env.service.Login(context.Background(), in)
	assertDomainErrorCode(t, err, domainerrors.CodeRateLimit)
}

func TestPasswordResetRequestRateLimited(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	in := usecase.PasswordResetRequestInput{Correo: correo, IP: "10.0.0.2", UserAgent: "ua"}
	for i := 0; i < 3; i++ {
		if err := env.service.PasswordResetRequest(context.Background(), in); err != nil {
			t.Fatalf("attempt %d: unexpected error: %v", i+1, err)
		}
	}

	err := env.service.PasswordResetRequest(context.Background(), in)
	assertDomainErrorCode(t, err, domainerrors.CodeRateLimit)
}

func TestPasswordResetRequestInactiveUserSendsNoEmail(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", false)

	err := env.service.PasswordResetRequest(context.Background(), usecase.PasswordResetRequestInput{
		Correo: correo, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("expected no error for an inactive user (avoid enumeration), got %v", err)
	}
	if len(env.mailer.Sent()) != 0 {
		t.Error("expected no email to be sent for an inactive user")
	}
}

func TestRefreshWithEmptyTokenFails(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Refresh(context.Background(), usecase.RefreshInput{IP: "127.0.0.1", UserAgent: "ua"})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestRefreshWithUnknownTokenFails(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: "token-that-was-never-issued", IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestRefreshWithExpiredTokenFails(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	userID := seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	// Seed an already-expired refresh token directly, bypassing issueSession
	// (which always sets a future expiry), to exercise the Expired() branch.
	const plainToken = "expired-refresh-token-for-test"
	_, err := env.pool.Exec(context.Background(),
		`INSERT INTO refresh_tokens (usuario_id, token_hash, expires_at)
		 VALUES ($1, $2, NOW() - INTERVAL '1 hour')`,
		userID, sha256Hex(plainToken),
	)
	if err != nil {
		t.Fatalf("seeding expired refresh token: %v", err)
	}

	_, err = env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: plainToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestMeWithUnknownUserFails(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Me(context.Background(), 9_999_999)
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestMeWithInactiveUserFails(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	userID := seedUsuario(t, env.pool, env.sedeID, correo, "vendedora", false)

	_, err := env.service.Me(context.Background(), userID)
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestPasswordResetConfirmWithExpiredTokenFails(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	userID := seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	const plainToken = "expired-reset-token-for-test"
	_, err := env.pool.Exec(context.Background(),
		`INSERT INTO password_resets (usuario_id, token_hash, expires_at)
		 VALUES ($1, $2, NOW() - INTERVAL '1 hour')`,
		userID, sha256Hex(plainToken),
	)
	if err != nil {
		t.Fatalf("seeding expired password reset: %v", err)
	}

	err = env.service.PasswordResetConfirm(context.Background(), usecase.PasswordResetConfirmInput{
		Token: plainToken, PasswordNueva: "otra-password-789", IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestPasswordResetConfirmWithUsedTokenFails(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	userID := seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	const plainToken = "already-used-reset-token-for-test"
	_, err := env.pool.Exec(context.Background(),
		`INSERT INTO password_resets (usuario_id, token_hash, expires_at, used_at)
		 VALUES ($1, $2, NOW() + INTERVAL '1 hour', NOW())`,
		userID, sha256Hex(plainToken),
	)
	if err != nil {
		t.Fatalf("seeding used password reset: %v", err)
	}

	err = env.service.PasswordResetConfirm(context.Background(), usecase.PasswordResetConfirmInput{
		Token: plainToken, PasswordNueva: "otra-password-789", IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestPasswordResetConfirmRejectsShortPassword(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.PasswordResetConfirm(context.Background(), usecase.PasswordResetConfirmInput{
		Token: "irrelevant-token", PasswordNueva: "short", IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeValidation)
}

func TestLogoutWithUnknownTokenSucceeds(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.Logout(context.Background(), usecase.LogoutInput{
		RefreshToken: "a-token-that-does-not-exist", IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("expected logout with an unknown token to succeed, got %v", err)
	}
}
