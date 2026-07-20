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

// Cerrar handles POST /api/v1/cuadres-caja/:id/cerrar. Admin-only, enforced
// by the router. A cuadre is immutable once cerrado — there is no reopen
// endpoint, so a second call against the same id returns 409.
func (h *Handler) Cerrar(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req CerrarCuadreRequest
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

	var valorTurno *decimal.Decimal
	if req.ValorTurno != nil {
		v, err := decimal.NewFromString(*req.ValorTurno)
		if err != nil {
			response.WriteError(w, r, domainerrors.NewValidation("Solicitud inválida", "valor_turno debe ser un número válido.", nil))
			return
		}
		valorTurno = &v
	}

	cuadre, err := h.service.Cerrar(r.Context(), usecase.CerrarInput{
		TargetID:      id,
		ValorTurno:    valorTurno,
		Observaciones: req.Observaciones,
		RequesterID:   requester.ID,
		IP:            middleware.IPFromContext(r.Context()),
		UserAgent:     middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteData(w, http.StatusOK, toCuadreResponse(*cuadre))
}
