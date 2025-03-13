package middleware

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "eduhub/server/internal/services/auth"
)

type KratosMiddleware struct {
    kratosService *auth.KratosService
}

func NewKratosMiddleware(kratosService *auth.KratosService) *KratosMiddleware {
    return &KratosMiddleware{
        kratosService: kratosService,
    }
}

func (m *KratosMiddleware) ValidateSession(next echo.HandlerFunc) echo.HandlerFunc {
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

        c.Set("identity", identity)
        return next(c)
    }
}

func (m *KratosMiddleware) RequireCollege(collegeID string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            identity, ok := c.Get("identity").(*auth.Identity)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "No identity found",
                })
            }

            if !m.kratosService.CheckCollegeAccess(identity, collegeID) {
                return c.JSON(http.StatusForbidden, map[string]string{
                    "error": "Access denied to this college",
                })
            }

            return next(c)
        }
    }
}

func (m *KratosMiddleware) RequireRole(roles ...string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            identity, ok := c.Get("identity").(*auth.Identity)
            if !ok {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "No identity found",
                })
            }

            for _, role := range roles {
                if m.kratosService.HasRole(identity, role) {
                    return next(c)
                }
            }

            return c.JSON(http.StatusForbidden, map[string]string{
                "error": "Insufficient role permissions",
            })
        }
    }
}