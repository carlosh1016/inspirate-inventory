package fragancias

import (
	"encoding/json"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
)

// Create handles POST /api/v1/fragancias.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateFraganciaRequest
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

	f, err := h.service.Create(r.Context(), usecase.CreateInput{
		SedeID:            requester.SedeID,
		NombreComercial:   req.NombreComercial,
		NombreAlternativo: req.NombreAlternativo,
		Genero:            req.Genero,
		NumeroGenero:      req.NumeroGenero,
		GramosMinimo:      req.GramosMinimo,
		RequesterID:       requester.ID,
		IP:                middleware.IPFromContext(r.Context()),
		UserAgent:         middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toFraganciaResponseFromGet(f))
}
