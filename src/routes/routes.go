package routes

import (
	"minitwit/src/handlers"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {
	// Timeline Routes
	app.GET("/", func(c echo.Context) error { return handlers.Timeline(c) })
	app.GET("/public", func(c echo.Context) error { return handlers.PublicTimeline(c) })
	app.GET("/:username", func(c echo.Context) error { return handlers.UserTimeline(c) })

	// User Follow Routes
	app.GET("/:username/follow", func(c echo.Context) error { return handlers.FollowUser(c) })
	app.GET("/:username/unfollow", func(c echo.Context) error { return handlers.UnfollowUser(c) })
	app.GET("/fllws/:username", func(c echo.Context) error { return handlers.Follow(c) })
	app.POST("/fllws/:username", func(c echo.Context) error { return handlers.Follow(c) })

	// Message Routes
	app.POST("/add_message", func(c echo.Context) error { return handlers.AddMessage(c) })
	app.GET("/msgs", func(c echo.Context) error { return handlers.Messages(c) })
	app.GET("/msgs/:username", func(c echo.Context) error { return handlers.MessagesPerUser(c) })
	app.POST("/msgs/:username", func(c echo.Context) error { return handlers.MessagesPerUser(c) })

	// Authentication Routes
	app.GET("/login", func(c echo.Context) error { return handlers.Login(c) })
	app.POST("/login", func(c echo.Context) error { return handlers.Login(c) })
	app.GET("/register", func(c echo.Context) error { return handlers.Register(c) })
	app.POST("/register", func(c echo.Context) error { return handlers.Register(c) })
	app.GET("/logout", func(c echo.Context) error { return handlers.Logout(c) })

	app.GET("/latest", func(c echo.Context) error { return handlers.GetLatest(c) })
	app.Static("/static", "static")

	// Prometheus metrics route
	app.GET("/metrics", echoprometheus.NewHandler())
}
