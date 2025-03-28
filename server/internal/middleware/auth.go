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

// AuthMiddleware uses AuthService to perform authentication (via Kratos) 
// and authorization (via Ory Keto) checks.
type AuthMiddleware struct {
	AuthService auth.AuthService
}

// NewAuthMiddleware now accepts an auth.AuthService instance,
// ensuring that the middleware has access to both authentication 
// (session validation) and authorization (permission checking) logic.
func NewAuthMiddleware(authSvc auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authSvc,
	}
}

// ValidateSession checks if the session token provided in the request 
// is valid. The AuthService.ValidateSession function should use Ory Kratos 
// to validate the session.
func (m *AuthMiddleware) ValidateSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionToken := c.Request().Header.Get("X-Session-Token")
		if sessionToken == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "No session token provided",
			})
		}

		identity, err := m.AuthService.ValidateSession(c.Request().Context(), sessionToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid session",
			})
		}

		// Store identity in context for later use by other middleware handlers.
		c.Set("identity", identity)
		return next(c)
	}
}

// RequireCollege ensures that the authenticated user belongs to the specified college.
// It extracts the collegeID from URL parameters and then calls AuthService.CheckCollegeAccess.
// Under a multitenant setup, this helps isolate college-specific resources.
func (m *AuthMiddleware) RequireCollege(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		collegeID := c.Param("collegeID")
		if collegeID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "College ID is required",
			})
		}

		identity, ok := c.Get("identity").(*auth.Identity)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
		}
		if !m.AuthService.CheckCollegeAccess(identity, collegeID) {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Access denied to this college",
			})
		}

		return next(c)
	}
}

// RequireRole ensures the user has at least one of the specified roles.
// The AuthService.HasRole method should return true for a given identity if it holds the role.
func (m *AuthMiddleware) RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identity, ok := c.Get("identity").(*auth.Identity)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}

			for _, role := range roles {
				if m.AuthService.HasRole(identity, role) {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Insufficient permissions",
			})
		}
	}
}

// RequirePermission ensures the user has the required permission on a given resource using Ory Keto.
// The AuthService.CheckPermission method should interface with the Ory Keto (via client-go)
// to check if the user (by identity.ID) has the specific action permitted on the resource.
// For example, resource could be a college identifier and action could be "view", "edit", etc.
func (m *AuthMiddleware) RequirePermission(resource, action string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identity, ok := c.Get("identity").(*auth.Identity)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}
			allowed, err := m.AuthService.CheckPermission(c.Request().Context(), identity, resource, action)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Error checking permissions",
				})
			}
			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Insufficient permissions",
				})
			}
			return next(c)
		}	}
}
