package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetNumber(c echo.Context) int {
    numStr := c.QueryParam("no")
	num := 100
	if numStr != "" {
		val, err := strconv.Atoi(numStr)
		if err == nil {
			num = val
		}
	}
	return num
}

