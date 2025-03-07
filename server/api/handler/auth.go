package handler

import (
	"eduhub/server/internal/helpers"
	"eduhub/server/internal/services/auth"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *auth.AuthService
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	OrgID     string `json:"orgId"`
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.JSONResponse(w, map[string]string{"error": "Invalid request"}, http.StatusBadRequest)
		return
	}

	// Validate organization
	if err := h.authService.ValidateOrganization(r.Context(), req.OrgID); err != nil {
		helpers.JSONResponse(w, map[string]string{"error": "Invalid organization"}, http.StatusBadRequest)
		return
	}

	// Register user
	err := h.authService.RegisterUser(r.Context(), auth.RegisterUserParams{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		OrgID:     req.OrgID,
	})
	if err != nil {
		helpers.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	helpers.JSONResponse(w, map[string]string{"message": "Registration successful"}, http.StatusCreated)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Redirect to Zitadel login page
	loginURL := h.authService.GetLoginURL()
	http.Redirect(w, r, loginURL, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		helpers.JSONResponse(w, map[string]string{"error": "No code provided"}, http.StatusBadRequest)
		return
	}

	token, err := h.authService.ExchangeCodeForToken(r.Context(), code)
	if err != nil {
		helpers.JSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Verify token and check organization
	claims, err := h.authService.VerifyToken(r.Context(), token.AccessToken())
	if err != nil {
		helpers.JSONResponse(w, map[string]string{"error": "Invalid token"}, http.StatusUnauthorized)
		return
	}

	helpers.JSONResponse(w, map[string]interface{}{
		"access_token": token.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   token.ExpiresIn,
		"user_id":      claims.Subject,
		"org_id":       claims.OrganizationID,
	}, http.StatusOK)
}
