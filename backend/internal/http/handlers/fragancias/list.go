package fragancias

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/fragancias"
)

// codigoFraganciaPattern recognizes the "M-014"/"F-3"/"f014" código shown in
// the fragancias table (see frontend lib/formatters.ts
// formatCodigoFragancia) — género initial, optional dash, digits.
var codigoFraganciaPattern = regexp.MustCompile(`(?i)^([mf])-?0*([0-9]+)$`)

// parseCodigoFragancia lets the free-text search box also match by código
// (not just nombre_comercial): the vendedora often knows a fragancia by its
// "F-014" label from the old Excel sheets, not by its commercial name.
func parseCodigoFragancia(q string) (genero string, numero int, ok bool) {
	m := codigoFraganciaPattern.FindStringSubmatch(strings.TrimSpace(q))
	if m == nil {
		return "", 0, false
	}
	n, err := strconv.Atoi(m[2])
	if err != nil {
		return "", 0, false
	}
	if strings.EqualFold(m[1], "m") {
		return "masculina", n, true
	}
	return "femenina", n, true
}

// List handles GET /api/v1/fragancias.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	requester, ok := middleware.UserFromContext(r.Context())
	if !ok {
		response.WriteError(w, r, domainerrors.NewUnauthorized("No autenticado", "Debes iniciar sesión para continuar."))
		return
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	numeroGenero, _ := strconv.Atoi(q.Get("numero_genero"))
	stockBajo, _ := strconv.ParseBool(q.Get("stock_bajo"))
	includeDeleted, _ := strconv.ParseBool(q.Get("include_deleted"))

	genero := q.Get("genero")
	texto := q.Get("q")
	if codigoGenero, codigoNumero, ok := parseCodigoFragancia(texto); ok {
		genero = codigoGenero
		numeroGenero = codigoNumero
		texto = ""
	}

	result, err := h.service.List(r.Context(), usecase.ListInput{
		Page:           page,
		PageSize:       pageSize,
		Sort:           q.Get("sort"),
		Q:              texto,
		SedeID:         requester.SedeID,
		Genero:         genero,
		NumeroGenero:   int32(numeroGenero),
		Activo:         q.Get("activo"),
		StockBajo:      stockBajo,
		IncludeDeleted: includeDeleted,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	items := make([]FraganciaResponse, len(result.Items))
	for i, f := range result.Items {
		items[i] = toFraganciaResponseFromList(f)
	}

	response.WriteList(w, http.StatusOK, items, result.Total, result.Page, result.PageSize)
}
