package ventas

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/shopspring/decimal"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	domainventas "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/ventas"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/ventas"
)

const idempotencyKeyHeader = "X-Idempotency-Key"

// Create handles POST /api/v1/ventas.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}

	var req CreateVentaRequest
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

	items, err := toCreateVentaItemInputs(req.Items)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	idempotencyKey, requestHash := idempotencyKeyFromRequest(r, body)

	out, err := h.service.CreateVenta(r.Context(), usecase.CreateVentaInput{
		SedeID:          requester.SedeID,
		UsuarioID:       requester.ID,
		MetodoPagoID:    req.MetodoPagoID,
		Items:           items,
		Observaciones:   req.Observaciones,
		IdempotencyKey:  idempotencyKey,
		RequestBodyHash: requestHash,
		IP:              middleware.IPFromContext(r.Context()),
		UserAgent:       middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	if out.FromIdempotency {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(int(out.CachedStatus))
		_, _ = w.Write(out.CachedBody)
		return
	}

	respBody, err := json.Marshal(map[string]any{"data": toVentaDetalladaResponse(*out.Venta, nil)})
	if err != nil {
		response.WriteError(w, r, domainerrors.NewInternal("Error inesperado", "No se pudo generar la respuesta.", err))
		return
	}

	if idempotencyKey != nil {
		if storeErr := h.service.StoreIdempotencyResponse(
			r.Context(), *idempotencyKey, requester.ID, requestHash, http.StatusCreated, respBody,
		); storeErr != nil {
			slog.ErrorContext(r.Context(), "failed to store idempotency response", "error", storeErr)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(respBody)
}

// idempotencyKeyFromRequest reads X-Idempotency-Key, treating an absent or
// empty header as "no idempotency" (nil key). When a non-empty key is
// present, it also hashes the raw request body — used to detect the same
// key being reused with a different body.
func idempotencyKeyFromRequest(r *http.Request, body []byte) (*string, string) {
	key := r.Header.Get(idempotencyKeyHeader)
	if key == "" {
		return nil, ""
	}
	sum := sha256.Sum256(body)
	return &key, hex.EncodeToString(sum[:])
}

func toCreateVentaItemInputs(items []CreateVentaItemRequest) ([]usecase.CreateVentaItemInput, error) {
	result := make([]usecase.CreateVentaItemInput, len(items))
	for i, it := range items {
		var gramos *decimal.Decimal
		if it.GramosFragancia != nil {
			g, err := decimal.NewFromString(*it.GramosFragancia)
			if err != nil {
				return nil, domainerrors.NewValidation(
					"Solicitud inválida",
					"gramos_fragancia debe ser un número válido.",
					nil,
				)
			}
			gramos = &g
		}
		result[i] = usecase.CreateVentaItemInput{
			TipoLinea:          domainventas.TipoLinea(it.TipoLinea),
			FraganciaID:        it.FraganciaID,
			VarianteEnvaseID:   it.VarianteEnvaseID,
			ProductoID:         it.ProductoID,
			FeromonaProductoID: it.FeromonaProductoID,
			GramosFragancia:    gramos,
			Cantidad:           it.Cantidad,
		}
	}
	return result, nil
}
