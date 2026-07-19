package movimientos

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /movimientos/* under r (typically the /api/v1 group). List,
// entrada-mercancia, traslado and danado are open to admin and vendedora;
// ajuste and correccion require admin.
func (h *Handler) Router(r chi.Router) {
	r.Route("/movimientos", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)
		r.Post("/entrada-mercancia", h.EntradaMercancia)
		r.Post("/traslado", h.Traslado)
		r.Post("/danado", h.Danado)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/ajuste", h.Ajuste)
			r.Post("/correccion", h.Correccion)
		})
	})
}
