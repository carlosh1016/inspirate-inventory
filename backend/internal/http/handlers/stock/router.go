package stock

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /stock/* under r (typically the /api/v1 group). Both
// endpoints are read-only and open to admin and vendedora.
func (h *Handler) Router(r chi.Router) {
	r.Route("/stock", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)
		r.Get("/alertas", h.Alertas)
	})
}
