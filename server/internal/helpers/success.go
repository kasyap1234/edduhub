package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Data any `json:"data"`
}

func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, SuccessResponse{Data: data})
}
