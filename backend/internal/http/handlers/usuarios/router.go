package usuarios

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /usuarios/* under r (typically the /api/v1 group). Every
// route requires an authenticated admin except PATCH /:id/password, which
// admin or the account owner may call (authorization checked in the
// handler and the usecase, not the middleware chain).
func (h *Handler) Router(r chi.Router) {
	r.Route("/usuarios", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))

		r.Patch("/{id}/password", h.UpdatePassword)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/", h.List)
			r.Get("/{id}", h.Get)
			r.Post("/", h.Create)
			r.Patch("/{id}", h.Update)
			r.Patch("/{id}/activar", h.Activate)
			r.Patch("/{id}/desactivar", h.Deactivate)
			r.Delete("/{id}", h.Delete)
		})
	})
}
