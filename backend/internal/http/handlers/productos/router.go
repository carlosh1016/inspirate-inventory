package productos

import (
	"github.com/go-chi/chi/v5"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
)

// Router mounts /productos/* under r (typically the /api/v1 group). List,
// get and update are open to admin and vendedora (the usecase restricts a
// vendedora's Update to stock_minimo only); create and delete require admin.
func (h *Handler) Router(r chi.Router) {
	r.Route("/productos", func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtManager))
		r.Use(middleware.RequireRole("admin", "vendedora"))

		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Patch("/{id}", h.Update)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Post("/", h.Create)
			r.Delete("/{id}", h.Delete)
		})
	})
}
