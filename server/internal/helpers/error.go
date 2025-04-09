package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	error any `json:"error"`
}

func Error(c echo.Context, error any) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		error: error,
	})
}
