package handler

import (
	"eduhub/server/internal/services/auth"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService 
}

func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// InitiateRegistration starts the registration flow
func (h *AuthHandler) InitiateRegistration(c echo.Context) error {
	flow, err := h.AuthService.InitiateRegistrationFlow(c.Request().Context())
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing flow ID",
		})
	}

	var req auth.RegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	identity, err := h.AuthService.CompleteRegistration(c.Request().Context(), flowID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, identity)
}

// HandleLogin processes login
func (h *AuthHandler) HandleLogin(c echo.Context) error {
	// Will be redirected to Kratos UI
	loginURL := fmt.Sprintf("%s/self-service/login/browser", h.AuthService.GetPublicURL())
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// HandleCallback processes the login callback
func (h *AuthHandler) HandleCallback(c echo.Context) error {
	sessionToken := c.Request().Header.Get("X-Session-Token")
	if sessionToken == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No session token provided",
		})
	}

	identity, err := h.AuthService.ValidateSession(c.Request().Context(), sessionToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid session",
		})
	}

	return c.JSON(http.StatusOK, identity)
}

			
			// InitiateRegistration starts the registration flow
// func (h *AuthHandler) InitiateRegistration(c echo.Context) error {
// 				flow, err := h.AuthService.InitiateRegistrationFlow(c.Request().Context())
// 				if err != nil {
// 					return c.JSON(http.StatusInternalServerError, map[string]string{
// 						"error": err.Error(),
// 					})
// 				}
// 				return c.JSON(http.StatusOK, flow)
// 			}
			
// 			// HandleRegistration processes the registration
// 			func (h *AuthHandler) HandleRegistration(c echo.Context) error {
// 				flowID := c.QueryParam("flow")
// 				if flowID == "" {
// 					return c.JSON(http.StatusBadRequest, map[string]string{
// 						"error": "Missing flow ID",
// 					})
// 				}
			
// 				var req auth.RegistrationRequest
// 				if err := c.Bind(&req); err != nil {
// 					return c.JSON(http.StatusBadRequest, map[string]string{
// 						"error": "Invalid request body",
// 					})
// 				}
			
// 				identity, err := h.KratosService.CompleteRegistration(c.Request().Context(), flowID, req)
// 				if err != nil {
// 					return c.JSON(http.StatusInternalServerError, map[string]string{
// 						"error": err.Error(),
// 					})
// 				}
			
// 				return c.JSON(http.StatusOK, identity)
// 			}
			
// 			// HandleLogin processes login
// 			func (h *AuthHandler) HandleLogin(c echo.Context) error {
// 				// Will be redirected to Kratos UI
// 				loginURL := fmt.Sprintf("%s/self-service/login/browser", h.KratosService.PublicURL)
// 				return c.Redirect(http.StatusTemporaryRedirect, loginURL)
// 			}
			
// 			// HandleCallback processes the login callback
// 			func (h *AuthHandler) HandleCallback(c echo.Context) error {
// 				sessionToken := c.Request().Header.Get("X-Session-Token")
// 				if sessionToken == "" {
// 					return c.JSON(http.StatusUnauthorized, map[string]string{
// 						"error": "No session token provided",
// 					})
// 				}
			
// 				identity, err := h.AuthService.ValidateSession(c.Request().Context(), sessionToken)
// 				if err != nil {
// 					return c.JSON(http.StatusUnauthorized, map[string]string{
// 						"error": "Invalid session",
// 					})
// 				}
			
// 				return c.JSON(http.StatusOK, identity)
// 			}
			
		
		

