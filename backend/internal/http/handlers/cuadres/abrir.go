package cuadres

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/cuadres"
)

// Abrir handles POST /api/v1/cuadres-caja/abrir. Admin-only, enforced by
// the router. The response includes a "warnings" array alongside "data" —
// a soft, non-blocking notice when a previous day's cuadre was left open.
func (h *Handler) Abrir(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req AbrirCuadreRequest
	if err := json.Unmarshal(body, &req); err != nil {
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

	var fondoBase *decimal.Decimal
	if req.FondoBase != nil {
		v, err := decimal.NewFromString(*req.FondoBase)
		if err != nil {
			response.WriteError(w, r, domainerrors.NewValidation("Solicitud inválida", "fondo_base debe ser un número válido.", nil))
			return
		}
		fondoBase = &v
	}

	out, err := h.service.Abrir(r.Context(), usecase.AbrirInput{
		SedeID:      requester.SedeID,
		FondoBase:   fondoBase,
		RequesterID: requester.ID,
		IP:          middleware.IPFromContext(r.Context()),
		UserAgent:   middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"data":     toCuadreResponse(*out.Cuadre),
		"warnings": toWarningResponses(out.Warnings),
	})
}
