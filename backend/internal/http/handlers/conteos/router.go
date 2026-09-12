package conteos

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /conteos-inventario/* under r (typically the /api/v1
// group). Everything is admin-only — solo la dueña crea, cuenta y cierra el
// conteo mensual; una vendedora no participa en este flujo.
func (h *Handler) Router(r chi.Router) {
	r.Route("/conteos-inventario", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin"))

		r.Get("/", h.List)
		r.Post("/", h.Crear)
		r.Get("/{id}", h.GetByID)
		r.Patch("/{id}/items/{item_id}", h.ActualizarItem)
		r.Post("/{id}/cerrar", h.Cerrar)
	})
}
