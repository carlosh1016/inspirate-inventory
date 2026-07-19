// Package response provides standardized JSON response helpers for HTTP
// handlers: successful payloads ({"data": ...}) and Problem Details errors.
package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// WriteJSON writes v as JSON with the given status code.
func WriteJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encode(w, body)
}

// WriteData wraps data in {"data": ...} and writes it as JSON.
func WriteData(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, map[string]any{"data": data})
}

// WriteNoContent writes an empty 204 response.
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Meta describes pagination metadata for list responses.
type Meta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// WriteList wraps items and pagination metadata in {"data": ..., "meta": ...}.
func WriteList(w http.ResponseWriter, status int, items any, total int64, page, pageSize int) {
	var totalPages int
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	WriteJSON(w, status, map[string]any{
		"data": items,
		"meta": Meta{Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages},
	})
}

func encode(w http.ResponseWriter, v any) {
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode json response", "error", err)
	}
}
