package utils

import "github.com/labstack/echo/v4"

func GetCsrfToken(c echo.Context) string {
    tok, ok := c.Get("csrf").(string)
    if !ok {
        return ""
    }
    return tok
}

func MapCSRFToContext(c echo.Context, data map[string]interface{}) map[string]interface{} {
    if data == nil {
        data = make(map[string]interface{})
    }
    data["CSRFToken"] = GetCsrfToken(c)
    return data
}