package modelosenvase

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/modelos_envase"
)

// Create handles POST /api/v1/modelos-envase.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateModeloEnvaseRequest
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

	sinVariantes := false
	if req.TieneVariantes != nil {
		sinVariantes = !*req.TieneVariantes
	}

	m, err := h.service.Create(r.Context(), usecase.CreateInput{
		Tipo:               req.Tipo,
		TamanoOz:           req.TamanoOz,
		EquivGramos:        req.EquivGramos,
		PrecioSolo:         req.PrecioSolo,
		PrecioConFragancia: req.PrecioConFragancia,
		PrecioRecarga:      req.PrecioRecarga,
		SinVariantes:       sinVariantes,
		SedeID:             requester.SedeID,
		RequesterID:        requester.ID,
		IP:                 middleware.IPFromContext(r.Context()),
		UserAgent:          middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toModeloEnvaseResponseFromGet(m))
}
