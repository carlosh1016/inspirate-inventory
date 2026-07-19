package modelosenvase

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
)

// Update handles PATCH /api/v1/modelos-envase/:id.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	var req UpdateModeloEnvaseRequest
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

	m, err := h.service.Update(r.Context(), usecase.UpdateInput{
		TargetID:           id,
		Tipo:               req.Tipo,
		TamanoOz:           req.TamanoOz,
		EquivGramos:        req.EquivGramos,
		PrecioSolo:         req.PrecioSolo,
		PrecioConFragancia: req.PrecioConFragancia,
		PrecioRecarga:      req.PrecioRecarga,
		RequesterID:        requester.ID,
		IP:                 middleware.IPFromContext(r.Context()),
		UserAgent:          middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toModeloEnvaseResponseFromGet(m))
}
