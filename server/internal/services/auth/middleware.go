package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
)

// VerifyToken middleware function that uses the auth service
func (s *AuthService) VerifyTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "No token provided",
			})
		}

		// Remove Bearer prefix if present
		if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		// Verify token using the auth service
		claims, err := s.VerifyToken(c.Request().Context(), token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
		}

		// Store claims in context
		c.Set("claims", claims)
		return next(c)
	}
}

// RequireOrganization middleware to check organization access
func (s *AuthService) RequireOrganization(orgID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get("claims").(*zitadel.TokenClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "No valid claims found",
				})
			}

			if claims.OrganizationID != orgID {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "User does not belong to this organization",
				})
			}

			return next(c)
		}
	}
}
