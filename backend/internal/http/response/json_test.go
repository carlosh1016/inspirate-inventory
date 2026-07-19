package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
)

func TestWriteData(t *testing.T) {
	w := httptest.NewRecorder()
	response.WriteData(w, http.StatusOK, map[string]string{"nombre": "Ana"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content-type: %q", ct)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %T", body["data"])
	}
	if data["nombre"] != "Ana" {
		t.Fatalf("expected nombre=Ana, got %v", data["nombre"])
	}
}

func TestWriteList(t *testing.T) {
	w := httptest.NewRecorder()
	items := []string{"a", "b", "c"}
	response.WriteList(w, http.StatusOK, items, 25, 2, 10)

	var body struct {
		Data []string      `json:"data"`
		Meta response.Meta `json:"meta"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(body.Data) != 3 {
		t.Fatalf("expected 3 items, got %d", len(body.Data))
	}
	if body.Meta.Total != 25 || body.Meta.Page != 2 || body.Meta.PageSize != 10 || body.Meta.TotalPages != 3 {
		t.Fatalf("unexpected meta: %+v", body.Meta)
	}
}

func TestWriteNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	response.WriteNoContent(w)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body, got %q", w.Body.String())
	}
}
