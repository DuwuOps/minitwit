package routes

import (
	"minitwit/src/handlers"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {
	// Timeline Routes
	app.GET("/", handlers.Timeline)
	app.GET("/public", handlers.PublicTimeline)
	app.GET("/:username", handlers.UserTimeline)

	// User Follow Routes
	app.GET("/:username/follow", handlers.FollowUser)
	app.GET("/:username/unfollow", handlers.UnfollowUser)
	app.GET("/fllws/:username", handlers.Follow)
	app.POST("/fllws/:username", handlers.Follow)

	// Message Routes
	app.POST("/add_message", handlers.AddMessage)
	app.GET("/msgs", handlers.Messages)
	app.GET("/msgs/:username", handlers.MessagesPerUser)
	app.POST("/msgs/:username", handlers.MessagesPerUser)

	// Authentication Routes
	app.GET("/login", handlers.Login)
	app.POST("/login", handlers.Login)
	app.GET("/register", handlers.Register)
	app.POST("/register", handlers.Register)
	app.GET("/logout", handlers.Logout)

	app.GET("/latest", func(c echo.Context) error { return handlers.GetLatest(c) })
	app.Static("/static", "static")

	// Prometheus metrics route
	app.GET("/metrics", echoprometheus.NewHandler())
}
