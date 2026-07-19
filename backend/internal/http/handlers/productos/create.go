package productos

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/productos"
)

// Create handles POST /api/v1/productos. Admin-only, enforced by the
// router.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.WriteError(w, r, err)
		return
	}

	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	var stockMinimo int32
	if req.StockMinimo != nil {
		stockMinimo = *req.StockMinimo
	}

	p, err := h.service.Create(r.Context(), usecase.CreateInput{
		SedeID:      requester.SedeID,
		Nombre:      req.Nombre,
		Categoria:   req.Categoria,
		Precio:      req.Precio,
		StockMinimo: stockMinimo,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toProductoResponseFromGet(p))
}
