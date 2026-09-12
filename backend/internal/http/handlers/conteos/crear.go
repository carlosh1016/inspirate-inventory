package conteos

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/conteos"
)

// Crear handles POST /api/v1/conteos-inventario. Admin-only, enforced by
// the router. Periodo defaults to the current month when omitted.
func (h *Handler) Crear(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req CrearConteoRequest
	if len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			response.WriteError(w, r, badRequestBodyErr())
			return
		}
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

	var periodo *time.Time
	if req.Periodo != nil {
		p, err := time.Parse("2006-01-02", *req.Periodo)
		if err != nil {
			response.WriteError(w, r, domainerrors.NewValidation("Solicitud inválida", "periodo debe tener el formato AAAA-MM-DD.", nil))
			return
		}
		periodo = &p
	}

	conteo, err := h.service.Crear(r.Context(), usecase.CrearInput{
		SedeID:      requester.SedeID,
		Periodo:     periodo,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusCreated, toConteoResponse(*conteo))
}
