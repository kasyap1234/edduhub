package middleware

import (
	"eduhub/server/internal/helpers"
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	RoleAdmin   = "admin"
	RoleFaculty = "faculty"
	RoleStudent = "student"

	identityContextKey  = "identity"
	collegeIDContextKey = "college_id"
	studentIDContextKey = "student_id"
	facultyIDContextKey = "faculty_id"
)

// AuthMiddleware uses AuthService to perform authentication (via Kratos)
// and authorization (via Ory Keto) checks.
type AuthMiddleware struct {
	AuthService auth.AuthService
	StudentRepo repository.StudentRepository
}

// NewAuthMiddleware now accepts an auth.AuthService instance,
// ensuring that the middleware has access to both authentication
// (session validation) and authorization (permission checking) logic.
func NewAuthMiddleware(authSvc auth.AuthService, studentRepo repository.StudentRepository) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authSvc,
		StudentRepo: studentRepo,
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
		c.Set(identityContextKey, identity)
		return next(c)
	}
}

// RequireCollege ensures that the authenticated user belongs to the specified college.
// It extracts the collegeID from URL parameters and then calls AuthService.CheckCollegeAccess.
// Under a multitenant setup, this helps isolate college-specific resources.
func (m *AuthMiddleware) RequireCollege(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		identity, ok := c.Get("identity").(*auth.Identity)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
		}
		userCollegeID := identity.Traits.College.ID
		c.Set("college_id", userCollegeID)

		return next(c)
	}
}

func (m *AuthMiddleware) LoadStudentProfile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		identity, ok := c.Get(identityContextKey).(*auth.Identity)
		if !ok || identity == nil {
			return helpers.Error(c, "Unauthorized")
		}
		if identity.Traits.Role == RoleStudent {
			student, err := m.StudentRepo.FindByKratosID(c.Request().Context(), identity.ID)
			if err != nil {
				return helpers.Error(c, "Unauthorized")
			}
			if student == nil {
				helpers.Error(c, "Not registered")
			}
			if !student.IsActive {
				return helpers.Error(c, "Inactive")
			}
			c.Set(studentIDContextKey, student.StudentID)
		}
		return next(c)
	}
}
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
		}
	}
}
