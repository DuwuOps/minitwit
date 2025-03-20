package routes

import (
	"database/sql"
	"minitwit/src/handlers"
	"minitwit/src/handlers/helpers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo, db *sql.DB) {
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
	app.GET("/login", handlers.Login)
	app.POST("/login", handlers.Login)
	app.GET("/register", handlers.Register)
	app.POST("/register", handlers.Register)
	app.GET("/logout", handlers.Logout)


	app.GET("/latest", func(c echo.Context) error { return helpers.GetLatest(c, db) })
	app.Static("/static", "static")
}
