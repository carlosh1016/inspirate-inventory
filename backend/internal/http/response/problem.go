package response

import (
	"errors"
	"log/slog"
	"net/http"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
)

const problemTypeBase = "https://api.inspirate.co/errors/"

// Problem is an RFC 7807 Problem Details body.
type Problem struct {
	Type     string              `json:"type"`
	Title    string              `json:"title"`
	Status   int                 `json:"status"`
	Detail   string              `json:"detail,omitempty"`
	Instance string              `json:"instance,omitempty"`
	Errors   map[string][]string `json:"errors,omitempty"`
}

// WriteError maps err to a Problem Details response. Non-DomainError values
// are logged with their real content and reduced to a generic 500 so
// internals never leak to the client.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var domainErr *domainerrors.DomainError
	if !errors.As(err, &domainErr) {
		slog.ErrorContext(r.Context(), "unhandled error", "error", err, "path", r.URL.Path)
		writeProblem(w, r, &domainerrors.DomainError{
			Code:   domainerrors.CodeInternal,
			Title:  "Error interno",
			Detail: "Ocurrió un error inesperado. Intenta de nuevo más tarde.",
		})
		return
	}

	if domainErr.Code == domainerrors.CodeInternal {
		slog.ErrorContext(r.Context(), "internal domain error", "error", domainErr.Wrapped, "path", r.URL.Path)
	}

	writeProblem(w, r, domainErr)
}

func writeProblem(w http.ResponseWriter, r *http.Request, domainErr *domainerrors.DomainError) {
	status := StatusForCode(domainErr.Code)
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)
	encode(w, Problem{
		Type:     problemTypeBase + string(domainErr.Code),
		Title:    domainErr.Title,
		Status:   status,
		Detail:   domainErr.Detail,
		Instance: r.URL.Path,
		Errors:   domainErr.Fields,
	})
}
