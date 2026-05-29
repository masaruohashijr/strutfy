package handlers

import (
	"encoding/json"
	"net/http"

	appauth "github.com/masaru/strutfy-api/internal/auth"

	"github.com/julienschmidt/httprouter"
)

type AuthHandler struct {
	auth *appauth.Service
}

func NewAuthHandler(auth *appauth.Service) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input appauth.RegisterInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, token, err := h.auth.Register(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user": map[string]any{
			"id":     user.ID,
			"org_id": user.OrganizationID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
		},
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   900,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input appauth.LoginInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, token, err := h.auth.Login(r.Context(), input)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": map[string]any{
			"id":     user.ID,
			"org_id": user.OrganizationID,
			"name":   user.Name,
			"email":  user.Email,
			"role":   user.Role,
		},
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   900,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}