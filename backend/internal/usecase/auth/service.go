// Package auth implements the login/refresh/logout/me/password-reset
// usecases: orchestration between domain/auth (pure rules) and the
// repository/platform layers. Every exported method returns either nil or a
// *domainerrors.DomainError — handlers never see raw repository errors.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"time"

	domainauth "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/auth"
	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/mailer"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/ratelimit"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/auditoria"
	passwordresets "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/password_resets"
	refreshtokens "github.com/carlosh1016/inspirate-inventory/backend/internal/repository/refresh_tokens"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/usuarios"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/repository/generated"
)

// Service groups every auth usecase behind one set of dependencies.
type Service struct {
	Usuarios       usuarios.Repository
	RefreshTokens  refreshtokens.Repository
	PasswordResets passwordresets.Repository
	Auditoria      auditoria.Repository
	JWT            jwt.Manager
	Mailer         mailer.Mailer
	LoginLimiter   ratelimit.Limiter
	ResetLimiter   ratelimit.Limiter

	JWTRefreshTTLAdmin     time.Duration
	JWTRefreshTTLVendedora time.Duration
	FrontendURL            string
}

// NewService builds a Service with all its dependencies.
func NewService(
	usuariosRepo usuarios.Repository,
	refreshTokensRepo refreshtokens.Repository,
	passwordResetsRepo passwordresets.Repository,
	auditoriaRepo auditoria.Repository,
	jwtManager jwt.Manager,
	mailerSvc mailer.Mailer,
	loginLimiter, resetLimiter ratelimit.Limiter,
	refreshTTLAdmin, refreshTTLVendedora time.Duration,
	frontendURL string,
) *Service {
	return &Service{
		Usuarios:               usuariosRepo,
		RefreshTokens:          refreshTokensRepo,
		PasswordResets:         passwordResetsRepo,
		Auditoria:              auditoriaRepo,
		JWT:                    jwtManager,
		Mailer:                 mailerSvc,
		LoginLimiter:           loginLimiter,
		ResetLimiter:           resetLimiter,
		JWTRefreshTTLAdmin:     refreshTTLAdmin,
		JWTRefreshTTLVendedora: refreshTTLVendedora,
		FrontendURL:            frontendURL,
	}
}

// issueSession generates a fresh access+refresh token pair for user and
// persists the refresh token's hash.
func (s *Service) issueSession(ctx context.Context, user generated.Usuario, ip, userAgent string) (domainauth.Session, error) {
	accessToken, accessExpiresAt, err := s.JWT.GenerateAccessToken(user.ID, string(user.Rol), user.SedeID)
	if err != nil {
		return domainauth.Session{}, internalErr(err)
	}

	plain, hash, err := generateOpaqueToken()
	if err != nil {
		return domainauth.Session{}, internalErr(err)
	}

	ttl := domainauth.RefreshTokenTTL(string(user.Rol), s.JWTRefreshTTLAdmin, s.JWTRefreshTTLVendedora)
	expiresAt := time.Now().Add(ttl)

	if _, err := s.RefreshTokens.Insert(ctx, user.ID, hash, ip, userAgent, expiresAt); err != nil {
		return domainauth.Session{}, internalErr(err)
	}

	return domainauth.Session{
		AccessToken:      accessToken,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     plain,
		RefreshExpiresAt: expiresAt,
	}, nil
}

// generateOpaqueToken returns a 32-byte random token, hex-encoded, plus the
// hex-encoded SHA-256 hash that gets persisted. The plaintext is only ever
// held in memory (cookie or email link), never stored.
func generateOpaqueToken() (plain, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	plain = hex.EncodeToString(buf)
	return plain, hashToken(plain), nil
}

func hashToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

// audit records an auditoria entry. Failures are logged, never returned —
// auditing must never block the primary action it's recording.
func (s *Service) audit(ctx context.Context, usuarioID *int64, accion, ip, userAgent string) {
	if err := s.Auditoria.Insert(ctx, auditoria.Entry{
		UsuarioID: usuarioID,
		Accion:    accion,
		IP:        ip,
		UserAgent: userAgent,
	}); err != nil {
		slog.ErrorContext(ctx, "failed to write auditoria entry", "accion", accion, "error", err)
	}
}

func internalErr(err error) error {
	return domainerrors.NewInternal("Error interno", "Ocurrió un error inesperado. Intenta de nuevo más tarde.", err)
}
