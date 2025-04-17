package helpers

import (
	"eduhub/server/internal/services/auth"

	"github.com/labstack/echo/v4"
)

func GetUserRole(c echo.Context)(string,error) {

	identity, ok :=c.Get("idenity").(*auth.Identity)
	if !ok {
		return "",echo.NewHTTPError(401,"unauthorized")
	}
	return identity.Traits.Role,nil 
}

