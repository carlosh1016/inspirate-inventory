package sesiones

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /sesiones-laborales/* under r (typically the /api/v1
// group). List is open to admin and vendedora (each scoped to her own
// sesiones, see List); entrada/salida are vendedora-only (a vendedora only
// clocks herself in/out); update and resumen are admin-only.
func (h *Handler) Router(r chi.Router) {
	r.Route("/sesiones-laborales", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("vendedora"))
			r.Post("/entrada", h.Entrada)
			r.Post("/salida", h.Salida)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/resumen", h.Resumen)
			r.Patch("/{id}", h.Update)
		})
	})
}
