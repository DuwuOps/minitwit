package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NotReqFromSimulator(c echo.Context) error {
	fromSimulator := c.Request().Header.Get("Authorization")
	if fromSimulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		data := map[string]any{
			"error_msg": "You are not authorized to use this resource!",
		}
		return c.JSON(http.StatusForbidden, data)
	}
	return nil
}
