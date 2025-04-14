package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotFoundResponse struct {
	Data any `json:"data"`
}


func NotFound(c echo.Context,data any)error {
	return c.JSON(http.StatusNotFound,NotFoundResponse{
		Data : data,
	}
)
}
