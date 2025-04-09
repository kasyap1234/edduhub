package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Data any `json:"data"`
}

func Error(c echo.Context, data any) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Data: data,
	})
}
