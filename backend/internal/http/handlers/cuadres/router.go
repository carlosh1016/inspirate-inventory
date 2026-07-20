package cuadres

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /cuadres-caja/* under r (typically the /api/v1 group).
// GetHoy and the pago/consignación create endpoints are open to admin and
// vendedora; everything else (List, GetByID, Abrir, Cerrar, the delete
// endpoints) is admin-only — only the admin opens/closes the day or
// removes an entry.
func (h *Handler) Router(r chi.Router) {
	r.Route("/cuadres-caja", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin", "vendedora"))
			r.Get("/hoy", h.GetHoy)
			r.Post("/{id}/pagos", h.CreatePago)
			r.Post("/{id}/consignaciones", h.CreateConsignacion)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/", h.List)
			r.Get("/{id}", h.GetByID)
			r.Post("/abrir", h.Abrir)
			r.Post("/{id}/cerrar", h.Cerrar)
			r.Delete("/{id}/pagos/{pago_id}", h.DeletePago)
			r.Delete("/{id}/consignaciones/{consignacion_id}", h.DeleteConsignacion)
		})
	})
}
