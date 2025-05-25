package middleware

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func Csrf() echo.MiddlewareFunc {
  return middleware.CSRFWithConfig(middleware.CSRFConfig{
        TokenLookup:	"header:X-CSRF-Token,form:_csrf",
        CookieName:		"csrf_token",
        CookiePath:		"/",
        CookieHTTPOnly:	true,
        CookieSecure:	false, // enable when serving HTTPS
    })
}