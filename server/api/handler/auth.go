package handler

import (
	"fmt"
	"net/http"

	"eduhub/server/internal/helpers"
	"eduhub/server/internal/services/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// InitiateRegistration starts the registration flow
func (h *AuthHandler) InitiateRegistration(c echo.Context) error {
	flow, err := h.authService.InitiateRegistrationFlow(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, flow)
}

// HandleRegistration processes the registration
func (h *AuthHandler) HandleRegistration(c echo.Context) error {
	flowID := c.QueryParam("flow")
	if flowID == "" {
		return helpers.Error(c, "empty flowID", 400)
	}

	var req auth.RegistrationRequest
	if err := c.Bind(&req); err != nil {
		return helpers.Error(c, "Invalid Registration Request", 400)
	}

	identity, err := h.authService.CompleteRegistration(c.Request().Context(), flowID, req)
	if err != nil {
		helpers.Error(c, "unable to complete registration", http.StatusNotFound)
	}

	return helpers.Success(c, identity, http.StatusOK)
}

// HandleLogin processes login
func (h *AuthHandler) HandleLogin(c echo.Context) error {
	// Will be redirected to Kratos UI
	loginURL := fmt.Sprintf("%s/self-service/login/browser", h.authService.GetPublicURL())
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// HandleCallback processes the login callback
func (h *AuthHandler) HandleCallback(c echo.Context) error {
	sessionToken := c.Request().Header.Get("X-Session-Token")
	if sessionToken == "" {
		return helpers.Error(c, "empty session token", 400)
	}

	identity, err := h.authService.ValidateSession(c.Request().Context(), sessionToken)
	if err != nil {
		return helpers.Error(c, "invalid identity", http.StatusInternalServerError)
	}

	return helpers.Success(c, identity, http.StatusOK)
}
