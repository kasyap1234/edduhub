package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotFoundResponse struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

func NotFound(c echo.Context, data interface{}, status int) error {
	return c.JSON(http.StatusNotFound, NotFoundResponse{
		Data:   data,
		Status: status,
	})
}
