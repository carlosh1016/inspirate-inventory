package auth

import (
	"net/http"
	"time"

	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/middleware"
	"github.com/carlosh1016/inspirate-inventory/backend/internal/http/response"
	usecase "github.com/carlosh1016/inspirate-inventory/backend/internal/usecase/auth"
)

type refreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	session, err := h.service.Refresh(r.Context(), usecase.RefreshInput{
		RefreshToken: refreshCookieFromRequest(r),
		IP:           middleware.IPFromContext(r.Context()),
		UserAgent:    middleware.UserAgentFromContext(r.Context()),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	h.setRefreshCookie(w, session.RefreshToken, session.RefreshExpiresAt)

	response.WriteData(w, http.StatusOK, refreshResponse{
		AccessToken: session.AccessToken,
		ExpiresIn:   int64(time.Until(session.AccessExpiresAt).Seconds()),
	})
}
