// Package response provee helpers para respuestas JSON estandarizadas.
package response

import (
	"encoding/json"
	"net/http"
)

// JSON escribe v como JSON en w con el status code indicado.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Error escribe un error estandarizado como JSON: {"error": {"message": "..."}}.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]any{
		"error": map[string]string{
			"message": message,
		},
	})
}
