package reportes

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /reportes/* under r (typically the /api/v1 group). Every report
// is admin-only — una vendedora nunca descarga reportes.
func (h *Handler) Router(r chi.Router) {
	r.Route("/reportes", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin"))

		r.Get("/ventas", h.Ventas)
		r.Get("/stock", h.Stock)
		r.Get("/movimientos", h.Movimientos)
		r.Get("/cuadres-caja", h.Cuadres)
		r.Get("/sesiones-laborales", h.Sesiones)
	})
}
