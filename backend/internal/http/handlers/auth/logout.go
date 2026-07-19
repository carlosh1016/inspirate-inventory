package auth

import (
	"net/http"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

// Logout handles POST /api/v1/auth/logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var userID *int64
	if user, ok := middleware.UserFromContext(r.Context()); ok {
		userID = &user.ID
	}

	err := h.service.Logout(r.Context(), usecase.LogoutInput{
		RefreshToken: refreshCookieFromRequest(r),
		IP:           middleware.IPFromContext(r.Context()),
		UserAgent:    middleware.UserAgentFromContext(r.Context()),
		UsuarioID:    userID,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	h.clearRefreshCookie(w)
	response.WriteNoContent(w)
}
