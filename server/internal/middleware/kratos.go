package middleware

import (
	"eduhub/server/internal/services/auth"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	RoleAdmin   = "admin"
	RoleFaculty = "faculty"
	RoleStudent = "student"
)

type AuthMiddleware struct {
	kratosService *auth.KratosService
}

func NewAuthMiddleware(kratosService *auth.KratosService) *AuthMiddleware {
	return &AuthMiddleware{
		kratosService: kratosService,
	}
}

// ValidateSession checks if the session is valid
func (m *AuthMiddleware) ValidateSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionToken := c.Request().Header.Get("X-Session-Token")
		if sessionToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "No session token provided",
			})
		}

		identity, err := m.kratosService.ValidateSession(c.Request().Context(), sessionToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid session",
			})
		}

		// Store identity in context
		c.Set("identity", identity)
		return next(c)
	}
}

// RequireCollege ensures user belongs to the specified college
func (m *AuthMiddleware) RequireCollege(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		collegeID := c.Param("collegeID")
		if collegeID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "College ID is required",
			})
		}

		identity := c.Get("identity").(*auth.Identity)
		if !m.kratosService.CheckCollegeAccess(identity, collegeID) {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied to this college",
			})
		}

		return next(c)
	}
}

// RequireRole ensures user has at least one of the specified roles
func (m *AuthMiddleware) RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identity := c.Get("identity").(*auth.Identity)

			for _, role := range roles {
				if m.kratosService.HasRole(identity, role) {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Insufficient permissions",
			})
		}
	}
}
