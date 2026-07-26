package metodospago

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /metodos-pago/* under r (typically the /api/v1 group).
// Reads (List, Get) are open to admin and vendedora — a vendedora needs the
// list to pick a método de pago when registering a venta. Writes (Create,
// Update, Delete) remain admin-only.
func (h *Handler) Router(r chi.Router) {
	r.Route("/metodos-pago", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin", "vendedora"))
			r.Get("/", h.List)
			r.Get("/{id}", h.Get)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", h.Create)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
}
