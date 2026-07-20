package ventas

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /ventas/* under r (typically the /api/v1 group). List, get,
// create and hoy/resumen are open to admin and vendedora (each usecase/
// handler enforces its own row-level scoping for a vendedora); update is
// admin-only.
func (h *Handler) Router(r chi.Router) {
	r.Route("/ventas", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)
		r.Get("/hoy/resumen", h.ResumenHoy)
		r.Get("/{id}", h.Get)
		r.Post("/", h.Create)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Patch("/{id}", h.Update)
		})
	})
}
