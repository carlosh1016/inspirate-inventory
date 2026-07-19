package auth

import (
	"encoding/json"
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

type passwordResetRequestRequest struct {
	Correo string `json:"correo" validate:"required,email"`
}

type passwordResetConfirmRequest struct {
	Token         string `json:"token" validate:"required"`
	PasswordNueva string `json:"password_nueva" validate:"required,min=8"`
}

// PasswordResetRequest handles POST /api/v1/auth/password-reset/request.
func (h *Handler) PasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	var req passwordResetRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.WriteError(w, r, err)
		return
	}

	err := h.service.PasswordResetRequest(r.Context(), usecase.PasswordResetRequestInput{
		Correo:    req.Correo,
		IP:        middleware.IPFromContext(r.Context()),
		UserAgent: middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteNoContent(w)
}

// PasswordResetConfirm handles POST /api/v1/auth/password-reset/confirm.
func (h *Handler) PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	var req passwordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.WriteError(w, r, err)
		return
	}

	err := h.service.PasswordResetConfirm(r.Context(), usecase.PasswordResetConfirmInput{
		Token:         req.Token,
		PasswordNueva: req.PasswordNueva,
		IP:            middleware.IPFromContext(r.Context()),
		UserAgent:     middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	response.WriteNoContent(w)
}
