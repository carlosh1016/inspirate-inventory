// Package auth contains the HTTP handlers for /api/v1/auth/*.
package auth

import (
	"net/http"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/jwt"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/platform/validator"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/v1/auth"
)

// Handler groups the auth HTTP handlers and their dependencies.
type Handler struct {
	service      *usecase.Service
	jwtManager   jwt.Manager
	validator    *validator.Validator
	isProduction bool
}

// NewHandler builds a Handler.
func NewHandler(service *usecase.Service, jwtManager jwt.Manager, v *validator.Validator, isProduction bool) *Handler {
	return &Handler{service: service, jwtManager: jwtManager, validator: v, isProduction: isProduction}
}

func (h *Handler) setRefreshCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isProduction,
		SameSite: http.SameSiteStrictMode,
	})
}

func refreshCookieFromRequest(r *http.Request) string {
	c, err := r.Cookie(refreshCookieName)
	if err != nil {
		return ""
	}
	return c.Value
}
