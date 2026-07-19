package middleware

import (
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

// RequireRole allows the request only if the authenticated user (set by
// Auth, earlier in the chain) has one of the given roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok || !allowed[user.Rol] {
				response.WriteError(w, r, domainerrors.NewForbidden(
					"Acceso denegado",
					"No tienes permiso para realizar esta acción.",
				))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
