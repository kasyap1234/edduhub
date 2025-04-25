package handler

import (
	"eduhub/server/internal/services/system"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type SystemHandler struct {
	systemService system.SystemService
}

// NewSystemHandler creates a new system handler
func NewSystemHandler(systemService system.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// @Summary Health check endpoint
// @Description Get server health status
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (s *SystemHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// SwaggerDocs serves the Swagger documentation
// @Summary API Documentation
// @Description Access the Swagger documentation
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /docs [get]
func (s *SystemHandler) SwaggerDocs(c echo.Context) error {
	return echoSwagger.WrapHandler(c)
}
