package metodospago

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /metodos-pago/* under r (typically the /api/v1 group).
// Every operation, including read, is admin-only.
func (h *Handler) Router(r chi.Router) {
	r.Route("/metodos-pago", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin"))

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Post("/", h.Create)
		r.Patch("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
