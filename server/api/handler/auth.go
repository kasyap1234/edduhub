package handler

import (
	"eduhub/server/internal/middleware"
	"eduhub/server/internal/models"
	"eduhub/server/internal/services"
	"eduhub/server/internal/services/auth"
	"eduhub/server/internal/services/student"
	"strconv"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService    auth.AuthService
	StudentService services.StudentService
}

func NewAuthHandler(authService auth.AuthService, studentService student.StudentService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		// initialise student service
		StudentService: studentService,
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
	ctx := c.Request().Context()
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

	identity, err := h.authService.CompleteRegistration(c.Request().Context(), flowID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	role := identity.Traits.Role
	collegeIDStr := identity.Traits.College.ID
	collegeID, err := strconv.Atoi(collegeIDStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	switch role {
	case middleware.RoleStudent:
		
		student := models.Student{
			KratosIdentityID: identity.ID,
			CollegeID:        collegeID,
			RollNo:           rollNo,
			IsActive:         true,
		}
		if err := h.StudentService.CreateStudent(ctx, &student); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}
	case middleware.RoleFaculty:
		// faculty based code
	case middleware.RoleAdmin:
		// admin based code
	}

	return c.JSON(http.StatusOK, identity)
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
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No session token provided",
		})
	}

	identity, err := h.authService.ValidateSession(c.Request().Context(), sessionToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid session",
		})
	}

	return c.JSON(http.StatusOK, identity)
}
