package auth

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /auth/* under r (typically the /api/v1 group).
func (h *Handler) Router(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/password-reset/request", h.PasswordResetRequest)
		r.Post("/password-reset/confirm", h.PasswordResetConfirm)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(h.jwtManager))
			r.Post("/logout", h.Logout)
			r.Get("/me", h.Me)
		})
	})
}
