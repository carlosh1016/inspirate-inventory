package movimientos

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/movimientos"
)

// Correccion handles POST /api/v1/movimientos/correccion.
func (h *Handler) Correccion(w http.ResponseWriter, r *http.Request) {
	var req CorreccionRequest
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

	m, err := h.service.Correccion(r.Context(), usecase.CorreccionInput{
		SedeID:        requester.SedeID,
		TipoItem:      req.TipoItem,
		ItemID:        req.ItemID,
		Ubicacion:     req.Ubicacion,
		CantidadNueva: req.CantidadNueva,
		Motivo:        req.Motivo,
		RequesterID:   requester.ID,
		IP:            middleware.IPFromContext(r.Context()),
		UserAgent:     middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	writeAjusteResult(w, m)
}
