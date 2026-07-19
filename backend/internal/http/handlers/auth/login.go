package auth

import (
	"encoding/json"
	"net/http"
	"time"

	domainerrors "github.com/carlosh1016/inspirate-inventory/backend/internal/domain/errors"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

type loginRequest struct {
	Correo   string `json:"correo" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken string          `json:"access_token"`
	ExpiresIn   int64           `json:"expires_in"`
	Usuario     usuarioResponse `json:"usuario"`
}

// Login handles POST /api/v1/auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, r, badRequestBodyErr())
		return
	}
	if err := h.validator.Validate(req); err != nil {
		response.WriteError(w, r, err)
		return
	}

	result, err := h.service.Login(r.Context(), usecase.LoginInput{
		Correo:    req.Correo,
		Password:  req.Password,
		IP:        middleware.IPFromContext(r.Context()),
		UserAgent: middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	h.setRefreshCookie(w, result.Session.RefreshToken, result.Session.RefreshExpiresAt)

	response.WriteData(w, http.StatusOK, loginResponse{
		AccessToken: result.Session.AccessToken,
		ExpiresIn:   int64(time.Until(result.Session.AccessExpiresAt).Seconds()),
		Usuario:     toUsuarioResponse(result.Usuario),
	})
}

func badRequestBodyErr() error {
	return domainerrors.NewValidation("Solicitud inválida", "El cuerpo de la solicitud no es un JSON válido.", nil)
}
