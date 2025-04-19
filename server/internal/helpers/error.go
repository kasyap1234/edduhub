package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message any `json:"error"`
	Status  int `json:"status"`
}

func Error(c echo.Context, error any, status int) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Message: error,
		Status:  status,
	})
}
