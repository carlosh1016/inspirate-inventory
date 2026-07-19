package middleware

import (
	"errors"
	"net/http"
	"strings"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
)

// Auth validates the Bearer access token and places the authenticated user
// in the request context for downstream handlers and RequireRole. Expired
// tokens get Detail="token_expired" so the frontend knows to call
// /auth/refresh instead of sending the user back to login.
func Auth(manager jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerToken(r)
			if !ok {
				response.WriteError(w, r, domainerrors.NewUnauthorized(
					"No autenticado",
					"Debes iniciar sesión para continuar.",
				))
				return
			}

			claims, err := manager.ParseAccessToken(token)
			if err != nil {
				detail := "El token de acceso no es válido."
				if errors.Is(err, jwt.ErrExpiredToken) {
					detail = "token_expired"
				}
				response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", detail))
				return
			}

			ctx := WithUser(r.Context(), AuthenticatedUser{
				ID:     claims.UserID,
				Rol:    claims.Rol,
				SedeID: claims.SedeID,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}
