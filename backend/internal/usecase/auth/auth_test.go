package auth_test

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

// uniqueCorreo derives a correo from the test name: correo is globally
// unique (uq_usuarios_correo has no sede_id), and every test in this
// package shares one Postgres container/pool.
func uniqueCorreo(t *testing.T) string {
	t.Helper()
	name := strings.ToLower(strings.NewReplacer("/", "-", " ", "-").Replace(t.Name()))
	return name + "@inspirate.co"
}

func TestLoginSuccess(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	result, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "test-agent",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Session.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if result.Session.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if result.Usuario.Correo != correo {
		t.Errorf("unexpected usuario: %+v", result.Usuario)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	_, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: "wrong-password", IP: "127.0.0.1", UserAgent: "test-agent",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestLoginUnknownCorreo(t *testing.T) {
	env := newTestEnv(t)

	_, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: uniqueCorreo(t), Password: "cualquiera123", IP: "127.0.0.1", UserAgent: "test-agent",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestLoginInactiveUser(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", false)

	_, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "test-agent",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeForbidden)
}

func TestRefreshRotatesToken(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	loginResult, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	newSession, err := env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if newSession.RefreshToken == loginResult.Session.RefreshToken {
		t.Error("expected refresh to rotate the token")
	}
	if newSession.AccessToken == "" {
		t.Error("expected a new access token")
	}

	// Reusing the just-rotated (now revoked) token must fail.
	_, err = env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestRefreshReuseRevokesAllSessions(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	loginResult, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	rotated, err := env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	// Reuse of the original (already-revoked) token: session is treated as
	// compromised, so it must revoke every token for the user...
	_, err = env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)

	// ...including the token that was legitimately rotated just above.
	_, err = env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: rotated.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestLogoutRevokesToken(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	loginResult, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	userID := loginResult.Usuario.ID
	err = env.service.Logout(context.Background(), usecase.LogoutInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua", UsuarioID: &userID,
	})
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	_, err = env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestLogoutWithoutTokenSucceeds(t *testing.T) {
	env := newTestEnv(t)

	if err := env.service.Logout(context.Background(), usecase.LogoutInput{IP: "127.0.0.1", UserAgent: "ua"}); err != nil {
		t.Fatalf("expected logout without a token to succeed, got %v", err)
	}
}

func TestPasswordResetRequestAndConfirmEndToEnd(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	err := env.service.PasswordResetRequest(context.Background(), usecase.PasswordResetRequestInput{
		Correo: correo, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	sent := env.mailer.Sent()
	if len(sent) != 1 {
		t.Fatalf("expected 1 email sent, got %d", len(sent))
	}
	token := extractTokenFromURL(t, sent[0].ResetURL)

	err = env.service.PasswordResetConfirm(context.Background(), usecase.PasswordResetConfirmInput{
		Token: token, PasswordNueva: "nueva-password-123", IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("confirm failed: %v", err)
	}

	_, err = env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)

	_, err = env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: "nueva-password-123", IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("expected login with the new password to succeed, got %v", err)
	}
}

func TestPasswordResetRequestUnknownCorreoStillSucceeds(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.PasswordResetRequest(context.Background(), usecase.PasswordResetRequestInput{
		Correo: uniqueCorreo(t), IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("expected no error for an unknown correo (avoid enumeration), got %v", err)
	}
	if len(env.mailer.Sent()) != 0 {
		t.Error("expected no email to be sent for an unknown correo")
	}
}

func TestPasswordResetConfirmInvalidatesSessions(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	seedUsuario(t, env.pool, env.sedeID, correo, "admin", true)

	loginResult, err := env.service.Login(context.Background(), usecase.LoginInput{
		Correo: correo, Password: testPassword, IP: "127.0.0.1", UserAgent: "ua",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if err := env.service.PasswordResetRequest(context.Background(), usecase.PasswordResetRequestInput{
		Correo: correo, IP: "127.0.0.1", UserAgent: "ua",
	}); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	sent := env.mailer.Sent()
	token := extractTokenFromURL(t, sent[len(sent)-1].ResetURL)

	if err := env.service.PasswordResetConfirm(context.Background(), usecase.PasswordResetConfirmInput{
		Token: token, PasswordNueva: "otra-password-456", IP: "127.0.0.1", UserAgent: "ua",
	}); err != nil {
		t.Fatalf("confirm failed: %v", err)
	}

	_, err = env.service.Refresh(context.Background(), usecase.RefreshInput{
		RefreshToken: loginResult.Session.RefreshToken, IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestPasswordResetConfirmWithInvalidTokenFails(t *testing.T) {
	env := newTestEnv(t)

	err := env.service.PasswordResetConfirm(context.Background(), usecase.PasswordResetConfirmInput{
		Token: "not-a-real-token", PasswordNueva: "cualquier-password-123", IP: "127.0.0.1", UserAgent: "ua",
	})
	assertDomainErrorCode(t, err, domainerrors.CodeUnauthorized)
}

func TestMeReturnsCurrentUsuario(t *testing.T) {
	env := newTestEnv(t)
	correo := uniqueCorreo(t)
	userID := seedUsuario(t, env.pool, env.sedeID, correo, "vendedora", true)

	usuario, err := env.service.Me(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if usuario.Correo != correo || string(usuario.Rol) != "vendedora" {
		t.Errorf("unexpected usuario: %+v", usuario)
	}
}

func assertDomainErrorCode(t *testing.T, err error, code domainerrors.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	var domainErr *domainerrors.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected *DomainError, got %T: %v", err, err)
	}
	if domainErr.Code != code {
		t.Fatalf("expected code %q, got %q (%v)", code, domainErr.Code, domainErr)
	}
}

func extractTokenFromURL(t *testing.T, resetURL string) string {
	t.Helper()
	u, err := url.Parse(resetURL)
	if err != nil {
		t.Fatalf("parsing reset url: %v", err)
	}
	token := u.Query().Get("token")
	if token == "" {
		t.Fatalf("no token found in reset url: %s", resetURL)
	}
	return token
}
