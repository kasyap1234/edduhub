package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Data   any `json:"data"`
	Status int `json:"status"`
}

func Success(c echo.Context, data any, status int) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Data:   data,
		Status: status,
	})
}
